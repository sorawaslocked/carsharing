package app

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	grpcserver "github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc"
	minioadapter "github.com/sorawaslocked/car-rental-car-service/internal/adapter/minio"
	natsadapter "github.com/sorawaslocked/car-rental-car-service/internal/adapter/nats"
	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/postgres"
	"github.com/sorawaslocked/car-rental-car-service/internal/config"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
	miniocfg "github.com/sorawaslocked/car-rental-car-service/internal/pkg/minio"
	natscfg "github.com/sorawaslocked/car-rental-car-service/internal/pkg/nats"
	postgrescfg "github.com/sorawaslocked/car-rental-car-service/internal/pkg/postgres"
	"github.com/sorawaslocked/car-rental-car-service/internal/service"
	"github.com/sorawaslocked/car-rental-car-service/internal/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log               *slog.Logger
	grpcServer        *grpcserver.Server
	telematicsService *service.TelematicsService
	natsConn          *nats.Conn
	db                *sql.DB
}

func New(
	cfg config.Config,
	logger *slog.Logger,
) (*App, error) {
	logger = logger.With(slog.String("appId", "car-service"))

	logger.Info("connecting to postgres database")
	db, err := postgrescfg.OpenDB(cfg.Postgres)
	if err != nil {
		logger.Error("connecting to postgres database", log.Err(err))
		return nil, err
	}

	logger.Info("connecting to minio")
	minioClient, err := miniocfg.NewClient(cfg.MinIO)
	if err != nil {
		logger.Error("connecting to minio", log.Err(err))
		return nil, err
	}
	objectStorage := minioadapter.NewObjectStorage(minioClient, cfg.MinIO.Bucket)

	logger.Info("connecting to NATS")
	nc, err := natscfg.Connect(cfg.NATS)
	if err != nil {
		logger.Error("connecting to NATS", log.Err(err))
		return nil, err
	}

	logger.Info("connecting to telematics stream service", slog.String("addr", cfg.TelematicsStream.Addr))
	telematicsConn, err := grpc.NewClient(cfg.TelematicsStream.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("connecting to telematics stream service", log.Err(err))
		return nil, err
	}

	validate := validator.New()
	if err = validation.RegisterCustomValidators(validate); err != nil {
		logger.Error("registering custom validators", log.Err(err))
		return nil, err
	}

	carModelRepo := postgres.NewCarModelRepository(db, logger)
	carRepo := postgres.NewCarRepository(db, logger)
	carInsuranceRepo := postgres.NewCarInsuranceRepository(db, logger)
	templateRepo := postgres.NewCarMaintenanceTemplateRepository(db, logger)
	recordRepo := postgres.NewCarMaintenanceRecordRepository(db, logger)
	serviceStateRepo := postgres.NewCarServiceStateRepository(db, logger)
	statusLogRepo := postgres.NewCarStatusLogRepository(db, logger)
	telematicsRepo := postgres.NewTelematicsRepository(db, logger)
	zoneRepo := postgres.NewZoneRepository(db, logger)

	natsPublisher := natsadapter.NewPublisher(nc)

	telematicsStreamClient := grpcserver.NewTelematicsStreamClient(telematicsConn, logger)
	telematicsService := service.NewTelematicsService(telematicsStreamClient, telematicsRepo, carRepo, logger)

	carModelService := service.NewCarModelService(carModelRepo, objectStorage, validate, logger)
	carService := service.NewCarService(carRepo, statusLogRepo, telematicsRepo, objectStorage, natsPublisher, validate, logger)
	carService.SetCarCreatedNotifier(telematicsService)

	carInsuranceService := service.NewCarInsuranceService(carInsuranceRepo, objectStorage, validate, logger)
	carMaintenanceService := service.NewCarMaintenanceService(templateRepo, recordRepo, serviceStateRepo, carRepo, carService, objectStorage, validate, logger)
	zoneService := service.NewZoneService(zoneRepo, validate, logger)

	natsSubscriber := natsadapter.NewSubscriber(nc, carService, logger)
	if err = natsSubscriber.Subscribe(); err != nil {
		logger.Error("subscribing to NATS events", log.Err(err))
		return nil, err
	}

	grpcSrv := grpcserver.NewServer(
		cfg.GRPC, logger,
		carModelService, carService, carInsuranceService, carMaintenanceService, zoneService,
		telematicsService,
		db, nc,
	)

	return &App{
		log:               logger,
		grpcServer:        grpcSrv,
		telematicsService: telematicsService,
		natsConn:          nc,
		db:                db,
	}, nil
}

func (a *App) stop() {
	a.telematicsService.Stop()
	a.grpcServer.Stop()
	if err := a.natsConn.Drain(); err != nil {
		a.log.Error("nats drain error", log.Err(err))
	}
	if err := a.db.Close(); err != nil {
		a.log.Error("postgres close error", log.Err(err))
	}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.grpcServer.MustRun()

	if err := a.telematicsService.Start(ctx); err != nil {
		a.log.Error("failed to start telematics service", log.Err(err))
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	s := <-shutdownCh

	a.log.Info("received system shutdown signal", slog.String("signal", s.String()))
	a.log.Info("stopping the application")
	cancel()
	a.stop()
	a.log.Info("graceful shutdown complete")
}
