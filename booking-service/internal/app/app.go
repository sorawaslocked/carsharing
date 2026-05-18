package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	natsgo "github.com/nats-io/nats.go"
	grpcadapter "github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc"
	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc/interceptor"
	natsadapter "github.com/sorawaslocked/car-rental-booking-service/internal/adapter/nats"
	"github.com/sorawaslocked/car-rental-booking-service/internal/adapter/postgres"
	"github.com/sorawaslocked/car-rental-booking-service/internal/config"
	pkglog "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/log"
	pkgnats "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/nats"
	pkgpostgres "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/postgres"
	"github.com/sorawaslocked/car-rental-booking-service/internal/service"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	grpcAddr   string
	db         *sql.DB
	natsConn   *natsgo.Conn
	cancel     context.CancelFunc
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	db, err := pkgpostgres.NewDB(cfg.PG)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	nc, err := pkgnats.NewConn(cfg.NATS)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("nats: %w", err)
	}

	bookingRepo := postgres.NewBookingRepo(log, db)
	ruleRepo := postgres.NewPricingRuleRepo(log, db)

	publisher := natsadapter.NewPublisher(log, nc)

	bookingSvc := service.NewBookingService(log, bookingRepo, ruleRepo, publisher)
	ruleSvc := service.NewPricingRuleService(log, ruleRepo)

	ctx, cancel := context.WithCancel(context.Background())

	consumer := natsadapter.NewConsumer(log, nc, bookingSvc)
	if err := consumer.Subscribe(ctx); err != nil {
		cancel()
		nc.Close()
		_ = db.Close()
		return nil, fmt.Errorf("nats consumer: %w", err)
	}

	go bookingSvc.StartExpiryWatcher(ctx)

	bookingHandler := handler.NewBookingHandler(log, bookingSvc)
	ruleHandler := handler.NewPricingRuleHandler(log, ruleSvc)
	healthHandler := handler.NewHealthHandler(log, db, nc)

	baseInterceptor := interceptor.NewBaseInterceptor()
	loggerInterceptor := interceptor.NewLoggerInterceptor(log)
	authInterceptor := interceptor.NewAuthInterceptor(log)

	grpcSrv := grpcadapter.NewServer(
		bookingHandler, ruleHandler, healthHandler,
		baseInterceptor, loggerInterceptor, authInterceptor,
	)

	return &App{
		log:        pkglog.WithComponent(log, "app"),
		grpcServer: grpcSrv,
		grpcAddr:   cfg.GRPC.Addr,
		db:         db,
		natsConn:   nc,
		cancel:     cancel,
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
	_ = a.db.Close()
	a.natsConn.Close()
}
