package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	grpcserver "github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc"
	grpcdocanalyzer "github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/client"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc/handler"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/mailer"
	minioadapter "github.com/sorawaslocked/car-rental-user-service/internal/adapter/minio"
	natsadapter "github.com/sorawaslocked/car-rental-user-service/internal/adapter/nats"
	natshandler "github.com/sorawaslocked/car-rental-user-service/internal/adapter/nats/handler"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/postgres"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/redis"
	"github.com/sorawaslocked/car-rental-user-service/internal/config"
	grpccfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/grpc"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	miniocfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/minio"
	natscfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/nats"
	postgrescfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/postgres"
	rediscfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/redis"
	validatecfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/validate"
	"github.com/sorawaslocked/car-rental-user-service/internal/service"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpcserver.Server
}

func New(
	_ context.Context,
	cfg config.Config,
	log *slog.Logger,
) (*App, error) {
	log = log.With(slog.String("appId", "user-service"))

	log.Info("connecting to postgres database")
	db, err := postgrescfg.OpenDB(cfg.Postgres)
	if err != nil {
		log.Error("connecting to postgres database", pkglog.Err(err))
		return nil, err
	}

	validate := validator.New()
	if err := validate.RegisterValidation("min_age", validatecfg.MinAge); err != nil {
		return nil, err
	}
	if err := validate.RegisterValidation("complex_password", validatecfg.ComplexPassword); err != nil {
		return nil, err
	}

	log.Info("connecting to nats")
	natsConn, err := natscfg.Connect(cfg.NATS)
	if err != nil {
		log.Error("connecting to nats", pkglog.Err(err))
		return nil, err
	}

	log.Info("connecting to document analyzer")
	analyzerConn, err := grpccfg.Connect(cfg.GRPCClient.Client.DocumentAnalyzerURL, cfg.GRPCClient.Client)
	if err != nil {
		log.Error("connecting to document analyzer", pkglog.Err(err))
		return nil, err
	}

	redisConn := rediscfg.Client(cfg.Redis)
	activationCodeCache := redis.NewActivationCodeRedisCache(redisConn)

	log.Info("connecting to minio")
	minioClient, err := miniocfg.NewClient(cfg.Minio)
	if err != nil {
		log.Error("connecting to minio", pkglog.Err(err))
		return nil, err
	}

	msMailer := mailer.New(log, cfg.Mailer)
	userRepo := postgres.NewUserRepository(log, db)
	docRepo := postgres.NewDocumentRepository(log, db)
	publisher := natsadapter.NewPublisher(log, natsConn)
	minioStorage := minioadapter.NewMinioObjectStorage(log, minioClient, cfg.Minio)
	analyzerClient := grpcdocanalyzer.NewDocumentAnalyzer(log, analyzerConn)

	userService := service.NewUserService(log, validate, userRepo, docRepo, minioStorage, analyzerClient, publisher, activationCodeCache, msMailer)

	log.Info("subscribing to nats document events")
	docHandler := natshandler.NewDocumentHandler(log, natsConn, userService)
	if err := docHandler.Subscribe(); err != nil {
		log.Error("subscribing to document analyzed events", pkglog.Err(err))
		return nil, err
	}

	healthHandler := handler.NewHealthHandler(log, map[string]handler.Pinger{
		"postgres":          postgres.NewChecker(db),
		"redis":             redis.NewChecker(redisConn),
		"nats":              natsadapter.NewChecker(natsConn),
		"minio":             minioadapter.NewChecker(minioClient, cfg.Minio.BucketName),
		"document-analyzer": grpcdocanalyzer.NewChecker(analyzerConn),
	})

	grpcServer := grpcserver.NewServer(cfg.GRPC, log, userService, healthHandler)

	return &App{
		log:        log,
		grpcServer: grpcServer,
	}, nil
}

func (a *App) stop() {
	a.grpcServer.Stop()
}

func (a *App) Run() {
	a.grpcServer.MustRun()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	s := <-shutdownCh

	a.log.Info("received system shutdown signal", slog.String("signal", s.String()))
	a.log.Info("stopping the application")
	a.stop()
	a.log.Info("graceful shutdown complete")
}
