package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	grpcserver "github.com/sorawaslocked/car-rental-car-service/internal/adapter/grpc"
	"github.com/sorawaslocked/car-rental-car-service/internal/adapter/postgres"
	"github.com/sorawaslocked/car-rental-car-service/internal/config"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
	postgrescfg "github.com/sorawaslocked/car-rental-car-service/internal/pkg/postgres"
	"github.com/sorawaslocked/car-rental-car-service/internal/service"
	"github.com/sorawaslocked/car-rental-car-service/internal/validation"

	"log/slog"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpcserver.Server
}

func New(
	cfg config.Config,
	logger *slog.Logger,
) (*App, error) {
	logger = logger.With(slog.String("appId", "user-service"))

	logger.Info("connecting to postgres database")
	db, err := postgrescfg.OpenDB(cfg.Postgres)
	if err != nil {
		logger.Error("connecting to postgres database", log.Err(err))

		return nil, err
	}

	validate := validator.New()
	err = validation.RegisterCustomValidators(validate)
	if err != nil {
		logger.Error("registering custom validators", log.Err(err))

		return nil, err
	}

	carModelRepo := postgres.NewCarModelRepository(db)
	carRepo := postgres.NewCarRepository(db)

	carModelService := service.NewCarModelService(carModelRepo, logger)
	carService := service.NewCarService(carRepo, logger)

	grpcServer := grpcserver.NewServer(cfg.GRPC, logger, carModelService, carService)

	return &App{
		log:        logger,
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
