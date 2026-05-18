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

	"github.com/go-playground/validator/v10"
	natsio "github.com/nats-io/nats.go"
	grpcserver "github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc"
	minioadapter "github.com/sorawaslocked/car-rental-car-service/internal/adapter/minio"
	natsadapter "github.com/sorawaslocked/car-rental-car-service/internal/adapter/nats"
	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/postgres"
	"github.com/sorawaslocked/car-rental-car-service/internal/config"
	pkggrpc "github.com/sorawaslocked/car-rental-car-service/internal/pkg/grpc"
	pkglog "github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
	pkgminio "github.com/sorawaslocked/car-rental-car-service/internal/pkg/minio"
	pkgnats "github.com/sorawaslocked/car-rental-car-service/internal/pkg/nats"
	pkgpostgres "github.com/sorawaslocked/car-rental-car-service/internal/pkg/postgres"
	"github.com/sorawaslocked/car-rental-car-service/internal/service"
	"github.com/sorawaslocked/car-rental-car-service/internal/validation"
	"google.golang.org/grpc"
)

type App struct {
	log               *slog.Logger
	grpcServer        *grpc.Server
	grpcAddr          string
	db                *sql.DB
	natsConn          *natsio.Conn
	telematicsConn    *grpc.ClientConn
	telematicsService *service.TelematicsService
}

func New(log *slog.Logger, cfg config.Config) (*App, error) {
	db, err := pkgpostgres.NewDB(cfg.PG)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	minioClient, err := pkgminio.NewClient(cfg.MinIO)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("minio: %w", err)
	}
	objectStorage := minioadapter.NewObjectStorage(minioClient, cfg.MinIO.Bucket)

	nc, err := pkgnats.NewConn(cfg.NATS)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("nats: %w", err)
	}

	telematicsConn, err := pkggrpc.NewClientConn(cfg.TelematicsStream.Addr)
	if err != nil {
		nc.Close()
		_ = db.Close()
		return nil, fmt.Errorf("telematics stream client: %w", err)
	}

	validate := validator.New()
	if err = validation.RegisterCustomValidators(validate); err != nil {
		_ = telematicsConn.Close()
		nc.Close()
		_ = db.Close()
		return nil, fmt.Errorf("register validators: %w", err)
	}

	carModelRepo := postgres.NewCarModelRepository(db, log)
	carRepo := postgres.NewCarRepository(db, log)
	carInsuranceRepo := postgres.NewCarInsuranceRepository(db, log)
	templateRepo := postgres.NewCarMaintenanceTemplateRepository(db, log)
	recordRepo := postgres.NewCarMaintenanceRecordRepository(db, log)
	serviceStateRepo := postgres.NewCarServiceStateRepository(db, log)
	statusLogRepo := postgres.NewCarStatusLogRepository(db, log)
	telematicsRepo := postgres.NewTelematicsRepository(db, log)
	zoneRepo := postgres.NewZoneRepository(db, log)

	natsPublisher := natsadapter.NewPublisher(nc)

	telematicsStreamClient := grpcserver.NewTelematicsStreamClient(telematicsConn, log)
	telematicsService := service.NewTelematicsService(telematicsStreamClient, telematicsRepo, carRepo, log)

	carModelService := service.NewCarModelService(carModelRepo, objectStorage, validate, log)
	carService := service.NewCarService(carModelRepo, carRepo, statusLogRepo, telematicsRepo, objectStorage, natsPublisher, validate, log)
	carService.SetCarCreatedNotifier(telematicsService)

	carInsuranceService := service.NewCarInsuranceService(carInsuranceRepo, objectStorage, validate, log)
	carMaintenanceService := service.NewCarMaintenanceService(templateRepo, recordRepo, serviceStateRepo, carRepo, carService, objectStorage, validate, log)
	zoneService := service.NewZoneService(zoneRepo, validate, log)

	natsSubscriber := natsadapter.NewSubscriber(nc, carService, log)
	if err = natsSubscriber.Subscribe(); err != nil {
		_ = telematicsConn.Close()
		nc.Close()
		_ = db.Close()
		return nil, fmt.Errorf("nats subscribe: %w", err)
	}

	grpcSrv := grpcserver.NewServer(
		log,
		carModelService, carService, carInsuranceService, carMaintenanceService, zoneService,
		telematicsService,
		db, nc,
	)

	return &App{
		log:               pkglog.WithComponent(log, "app"),
		grpcServer:        grpcSrv.GRPCServer(),
		grpcAddr:          cfg.GRPC.Addr,
		db:                db,
		natsConn:          nc,
		telematicsConn:    telematicsConn,
		telematicsService: telematicsService,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	if err := a.telematicsService.Start(ctx); err != nil {
		cancel()
		return fmt.Errorf("telematics service start: %w", err)
	}

	lis, err := net.Listen("tcp", a.grpcAddr)
	if err != nil {
		cancel()
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
		cancel()
		return fmt.Errorf("gRPC server stopped unexpectedly: %w", err)
	}

	cancel()
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

	a.telematicsService.Stop()
	_ = a.db.Close()
	a.natsConn.Close()
	_ = a.telematicsConn.Close()
}
