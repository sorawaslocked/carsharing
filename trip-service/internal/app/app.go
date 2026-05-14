package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"

	natsio "github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	grpcserver "github.com/sorawaslocked/car-rental-trip-service/internal/adapter/grpc"
	"github.com/sorawaslocked/car-rental-trip-service/internal/adapter/grpc/client"
	"github.com/sorawaslocked/car-rental-trip-service/internal/adapter/grpc/handler"
	natspub "github.com/sorawaslocked/car-rental-trip-service/internal/adapter/nats"
	"github.com/sorawaslocked/car-rental-trip-service/internal/adapter/postgres"
	"github.com/sorawaslocked/car-rental-trip-service/internal/config"
	pkggrpc "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/grpc"
	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
	pkgnats "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/nats"
	pkgpostgres "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/postgres"
	"github.com/sorawaslocked/car-rental-trip-service/internal/service"
)

type App struct {
	log         *slog.Logger
	grpcServer  *grpc.Server
	grpcPort    int
	db          *sql.DB
	natsConn    *natsio.Conn
	carConn     *grpc.ClientConn
	streamConn  *grpc.ClientConn
	bookingConn *grpc.ClientConn
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	db, err := pkgpostgres.NewDB(cfg.PG.DSN)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	nc, err := pkgnats.NewConn(cfg.NATS.URL)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("nats: %w", err)
	}

	carConn, err := pkggrpc.NewClientConn(cfg.CarService.Addr)
	if err != nil {
		_ = db.Close()
		nc.Close()
		return nil, fmt.Errorf("car service client: %w", err)
	}

	streamConn, err := pkggrpc.NewClientConn(cfg.CarStreamService.Addr)
	if err != nil {
		_ = db.Close()
		nc.Close()
		_ = carConn.Close()
		return nil, fmt.Errorf("car stream service client: %w", err)
	}

	bookingConn, err := pkggrpc.NewClientConn(cfg.BookingService.Addr)
	if err != nil {
		_ = db.Close()
		nc.Close()
		_ = carConn.Close()
		_ = streamConn.Close()
		return nil, fmt.Errorf("booking service client: %w", err)
	}

	tripRepo := postgres.NewTripRepo(log, db)
	summaryRepo := postgres.NewTripSummaryRepo(log, db)
	statusRepo := postgres.NewTripStatusReadingRepo(log, db)

	bookingClient := client.NewBookingClient(log, bookingConn)
	telematicsClient := client.NewTelematicsClient(log, carConn, streamConn)
	publisher := natspub.NewPublisher(log, nc)

	tripSvc := service.NewTripService(log, tripRepo, summaryRepo, statusRepo, bookingClient, telematicsClient, publisher)

	tripHandler := handler.NewTripHandler(log, tripSvc)
	streamHandler := handler.NewTripStreamHandler(log, tripSvc)
	healthHandler := handler.NewHealthHandler(log)

	srv := grpcserver.NewServer(log, tripHandler, streamHandler, healthHandler)

	return &App{
		log:         pkglog.WithComponent(log, "app"),
		grpcServer:  srv,
		grpcPort:    cfg.GRPC.Port,
		db:          db,
		natsConn:    nc,
		carConn:     carConn,
		streamConn:  streamConn,
		bookingConn: bookingConn,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	a.log.Info("gRPC server listening", slog.Int("port", a.grpcPort))
	return a.grpcServer.Serve(lis)
}

func (a *App) Stop() {
	a.log.Info("shutting down")
	a.grpcServer.GracefulStop()
	_ = a.db.Close()
	a.natsConn.Close()
	_ = a.carConn.Close()
	_ = a.streamConn.Close()
	_ = a.bookingConn.Close()
}
