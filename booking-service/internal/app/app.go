package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcadapter "carsharing/booking-service/internal/adapter/grpc"
	grpcclient "carsharing/booking-service/internal/adapter/grpc/client"
	"carsharing/booking-service/internal/adapter/grpc/handler"
	"carsharing/booking-service/internal/adapter/grpc/interceptor"
	natsadapter "carsharing/booking-service/internal/adapter/nats"
	natspub "carsharing/booking-service/internal/adapter/nats/publisher"
	natssub "carsharing/booking-service/internal/adapter/nats/subscriber"
	postgresadapter "carsharing/booking-service/internal/adapter/postgres"
	"carsharing/booking-service/internal/config"
	"carsharing/booking-service/internal/service"
	"carsharing/booking-service/internal/validation"
	pkggrpc "carsharing/shared/pkg/grpc"
	pkglog "carsharing/shared/pkg/log"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	grpcAddr   string
	closer     closer
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	var cl closer

	pool, err := pkgpostgres.NewPool(log, cfg.PG)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	cl.add(pool.Close)

	publisherConn, err := pkgnats.NewPublisher(log, cfg.NATSPublisher)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats publisher: %w", err)
	}
	cl.add(publisherConn.Close)

	subscriberConn, err := pkgnats.NewSubscriber(log, cfg.NATSSubscriber)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats subscriber: %w", err)
	}
	cl.add(subscriberConn.Close)

	validate := validator.New()
	if err := validation.RegisterCustomValidators(validate, log); err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("validator: %w", err)
	}

	baseClientInterceptor := interceptor.NewClientBaseInterceptor()
	carServiceConn, err := pkggrpc.NewClientConn(
		log,
		cfg.CarService,
		grpc.WithChainUnaryInterceptor(baseClientInterceptor.Unary),
		grpc.WithChainStreamInterceptor(),
	)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("car service client: %w", err)
	}
	cl.add(func() { _ = carServiceConn.Close() })

	bookingRepo := postgresadapter.NewBookingRepository(log, pool)
	ruleRepo := postgresadapter.NewPricingRuleRepository(log, pool)

	publisher := natspub.NewBookingPublisher(log, publisherConn)

	carChecker := grpcclient.NewCarChecker(log, carServiceConn)
	carModelChecker := grpcclient.NewCarModelChecker(log, carServiceConn)

	bookingSvc := service.NewBookingService(log, validate, bookingRepo, publisher, carChecker)
	ruleSvc := service.NewPricingRuleService(log, validate, ruleRepo, carModelChecker)

	ctx, cancel := context.WithCancel(context.Background())
	cl.add(cancel)

	tripSubscriber := natssub.NewTripSubscriber(log, subscriberConn, bookingSvc)
	if err := tripSubscriber.Subscribe(); err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats subscriber: %w", err)
	}

	go bookingSvc.StartExpiryWatcher(ctx)

	bookingHandler := handler.NewBookingHandler(log, bookingSvc)
	ruleHandler := handler.NewPricingRuleHandler(log, ruleSvc)
	healthHandler := handler.NewHealthHandler(log, map[string]handler.Pinger{
		"postgres":    pool,
		"nats":        natsadapter.NewPinger(log, publisherConn),
		"car-service": grpcclient.NewPinger(log, carServiceConn),
	}, cfg.Version)

	grpcSrv, err := grpcadapter.NewServer(log, cfg.GRPC, bookingHandler, ruleHandler, healthHandler)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("grpc server: %w", err)
	}

	return &App{
		log:        pkglog.WithComponent(log, "app"),
		grpcServer: grpcSrv,
		grpcAddr:   fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port),
		closer:     cl,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", a.grpcAddr)
	if err != nil {
		return fmt.Errorf("gRPC listen: %w", err)
	}
	a.log.Info("gRPC server listening", slog.String("addr", a.grpcAddr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.grpcServer.Serve(lis)
	}()

	select {
	case sig := <-quit:
		a.log.Info("received signal, shutting down", slog.String("signal", sig.String()))
	case err := <-errCh:
		return fmt.Errorf("gRPC server stopped unexpectedly: %w", err)
	}

	a.stop()

	return nil
}

func (a *App) stop() {
	a.log.Info("shutting down")

	stopped := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		a.log.Info("gRPC server stopped gracefully")
	case <-time.After(15 * time.Second):
		a.log.Warn("graceful stop timed out, forcing stop")
		a.grpcServer.Stop()
	}

	a.closer.closeAll()
}
