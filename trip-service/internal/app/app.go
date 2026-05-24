package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"

	pkggrpc "carsharing/shared/pkg/grpc"
	pkglog "carsharing/shared/pkg/log"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"
	grpcserver "carsharing/trip-service/internal/adapter/grpc"
	"carsharing/trip-service/internal/adapter/grpc/client"
	"carsharing/trip-service/internal/adapter/grpc/handler"
	"carsharing/trip-service/internal/adapter/grpc/interceptor"
	natspub "carsharing/trip-service/internal/adapter/nats"
	"carsharing/trip-service/internal/adapter/postgres"
	"carsharing/trip-service/internal/config"
	"carsharing/trip-service/internal/service"
	"carsharing/trip-service/internal/validation"
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

	validate := validator.New()
	if err := validation.RegisterCustomValidators(validate, log); err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("validator: %w", err)
	}

	nc, err := pkgnats.NewPublisher(log, cfg.NATS)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("nats: %w", err)
	}
	cl.add(nc.Close)

	unaryInterceptors := grpc.WithChainUnaryInterceptor(interceptor.MetadataForwardingUnaryInterceptor)
	streamInterceptors := grpc.WithChainStreamInterceptor(interceptor.MetadataForwardingStreamInterceptor)

	carConn, err := pkggrpc.NewClientConn(log, cfg.CarService, unaryInterceptors, streamInterceptors)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("car service client: %w", err)
	}
	cl.add(func() { _ = carConn.Close() })

	streamConn, err := pkggrpc.NewClientConn(log, cfg.CarStreamService, unaryInterceptors, streamInterceptors)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("car stream service client: %w", err)
	}
	cl.add(func() { _ = streamConn.Close() })

	bookingConn, err := pkggrpc.NewClientConn(log, cfg.BookingService, unaryInterceptors, streamInterceptors)
	if err != nil {
		cl.closeAll()
		return nil, fmt.Errorf("booking service client: %w", err)
	}
	cl.add(func() { _ = bookingConn.Close() })

	tripRepo := postgres.NewTripRepo(log, pool)
	summaryRepo := postgres.NewTripSummaryRepo(log, pool)
	statusRepo := postgres.NewTripStatusReadingRepo(log, pool)

	bookingClient := client.NewBookingClient(log, bookingConn)
	telematicsClient := client.NewTelematicsClient(log, carConn, streamConn)
	publisher := natspub.NewPublisher(log, nc)

	tripSvc := service.NewTripService(log, validate, tripRepo, summaryRepo, statusRepo, bookingClient, telematicsClient, publisher)

	tripHandler := handler.NewTripHandler(log, tripSvc)
	streamHandler := handler.NewTripStreamHandler(log, tripSvc)
	healthHandler := handler.NewHealthHandler(log, pool, nc, carConn, streamConn, bookingConn)

	srv := grpcserver.NewServer(log, tripHandler, streamHandler, healthHandler)

	return &App{
		log:        pkglog.WithComponent(log, "app"),
		grpcServer: srv,
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
