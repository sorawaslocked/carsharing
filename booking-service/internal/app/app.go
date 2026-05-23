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
	"carsharing/booking-service/internal/adapter/grpc/handler"
	"carsharing/booking-service/internal/adapter/grpc/interceptor"
	natsadapter "carsharing/booking-service/internal/adapter/nats"
	postgresadapter "carsharing/booking-service/internal/adapter/postgres"
	"carsharing/booking-service/internal/config"
	"carsharing/booking-service/internal/service"
	pkglog "carsharing/shared/pkg/log"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	natsgo "github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type App struct {
	log            *slog.Logger
	grpcServer     *grpc.Server
	grpcAddr       string
	pool           *pgxpool.Pool
	publisherConn  *natsgo.Conn
	subscriberConn *natsgo.Conn
	cancel         context.CancelFunc
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	pool, err := pkgpostgres.NewPool(log, cfg.PG)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	publisherConn, err := pkgnats.NewPublisher(log, cfg.NATSPublisher)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("nats publisher: %w", err)
	}

	subscriberConn, err := pkgnats.NewSubscriber(log, cfg.NATSSubscriber)
	if err != nil {
		pool.Close()
		publisherConn.Close()
		return nil, fmt.Errorf("nats subscriber: %w", err)
	}

	bookingRepo := postgresadapter.NewBookingRepo(log, pool)
	ruleRepo := postgresadapter.NewPricingRuleRepo(log, pool)

	publisher := natsadapter.NewPublisher(log, publisherConn)

	bookingSvc := service.NewBookingService(log, bookingRepo, ruleRepo, publisher)
	ruleSvc := service.NewPricingRuleService(log, ruleRepo)

	ctx, cancel := context.WithCancel(context.Background())

	consumer := natsadapter.NewConsumer(log, subscriberConn, bookingSvc)
	if err := consumer.Subscribe(ctx); err != nil {
		cancel()
		subscriberConn.Close()
		publisherConn.Close()
		pool.Close()
		return nil, fmt.Errorf("nats consumer: %w", err)
	}

	go bookingSvc.StartExpiryWatcher(ctx)

	bookingHandler := handler.NewBookingHandler(log, bookingSvc)
	ruleHandler := handler.NewPricingRuleHandler(log, ruleSvc)
	healthHandler := handler.NewHealthHandler(log, pool, publisherConn)

	baseInterceptor := interceptor.NewBaseInterceptor()
	loggerInterceptor := interceptor.NewLoggerInterceptor(log)
	authInterceptor := interceptor.NewAuthInterceptor(log)

	grpcSrv, err := grpcadapter.NewServer(
		log, cfg.GRPC,
		bookingHandler, ruleHandler, healthHandler,
		baseInterceptor, loggerInterceptor, authInterceptor,
	)
	if err != nil {
		cancel()
		subscriberConn.Close()
		publisherConn.Close()
		pool.Close()
		return nil, fmt.Errorf("grpc server: %w", err)
	}

	return &App{
		log:            pkglog.WithComponent(log, "app"),
		grpcServer:     grpcSrv,
		grpcAddr:       fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port),
		pool:           pool,
		publisherConn:  publisherConn,
		subscriberConn: subscriberConn,
		cancel:         cancel,
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
		a.log.Info("grpc server stopped gracefully")
	case <-time.After(15 * time.Second):
		a.log.Warn("graceful stop timed out, forcing stop")
		a.grpcServer.Stop()
	}

	a.cancel()
	a.pool.Close()
	a.publisherConn.Close()
	a.subscriberConn.Close()
}
