package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	natsio "github.com/nats-io/nats.go"
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
)

type App struct {
	log         *slog.Logger
	grpcServer  *grpc.Server
	grpcAddr    string
	pool        *pgxpool.Pool
	natsConn    *natsio.Conn
	carConn     *grpc.ClientConn
	streamConn  *grpc.ClientConn
	bookingConn *grpc.ClientConn
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	pool, err := pkgpostgres.NewPool(log, cfg.PG)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	nc, err := pkgnats.NewPublisher(log, cfg.NATS)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("nats: %w", err)
	}

	unaryInterceptors := grpc.WithChainUnaryInterceptor(interceptor.MetadataForwardingUnaryInterceptor)
	streamInterceptors := grpc.WithChainStreamInterceptor(interceptor.MetadataForwardingStreamInterceptor)

	carConn, err := pkggrpc.NewClientConn(log, cfg.CarService, unaryInterceptors, streamInterceptors)
	if err != nil {
		pool.Close()
		nc.Close()
		return nil, fmt.Errorf("car service client: %w", err)
	}

	streamConn, err := pkggrpc.NewClientConn(log, cfg.CarStreamService, unaryInterceptors, streamInterceptors)
	if err != nil {
		pool.Close()
		nc.Close()
		_ = carConn.Close()
		return nil, fmt.Errorf("car stream service client: %w", err)
	}

	bookingConn, err := pkggrpc.NewClientConn(log, cfg.BookingService, unaryInterceptors, streamInterceptors)
	if err != nil {
		pool.Close()
		nc.Close()
		_ = carConn.Close()
		_ = streamConn.Close()
		return nil, fmt.Errorf("booking service client: %w", err)
	}

	tripRepo := postgres.NewTripRepo(log, pool)
	summaryRepo := postgres.NewTripSummaryRepo(log, pool)
	statusRepo := postgres.NewTripStatusReadingRepo(log, pool)

	bookingClient := client.NewBookingClient(log, bookingConn)
	telematicsClient := client.NewTelematicsClient(log, carConn, streamConn)
	publisher := natspub.NewPublisher(log, nc)

	tripSvc := service.NewTripService(log, tripRepo, summaryRepo, statusRepo, bookingClient, telematicsClient, publisher)

	tripHandler := handler.NewTripHandler(log, tripSvc)
	streamHandler := handler.NewTripStreamHandler(log, tripSvc)
	healthHandler := handler.NewHealthHandler(log, pool, nc, carConn, streamConn, bookingConn)

	srv := grpcserver.NewServer(log, tripHandler, streamHandler, healthHandler)

	return &App{
		log:         pkglog.WithComponent(log, "app"),
		grpcServer:  srv,
		grpcAddr:    fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port),
		pool:        pool,
		natsConn:    nc,
		carConn:     carConn,
		streamConn:  streamConn,
		bookingConn: bookingConn,
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

	a.pool.Close()
	a.natsConn.Close()
	_ = a.carConn.Close()
	_ = a.streamConn.Close()
	_ = a.bookingConn.Close()
}
