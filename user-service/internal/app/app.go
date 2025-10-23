package app

import (
	grpcserver "car-rental-user-service/internal/adapter/grpc"
	"car-rental-user-service/internal/adapter/postgres"
	"car-rental-user-service/internal/config"
	"car-rental-user-service/internal/pkg/jwt"
	"car-rental-user-service/internal/pkg/logger"
	postgrescfg "car-rental-user-service/internal/pkg/postgres"
	validatecfg "car-rental-user-service/internal/pkg/validate"
	"car-rental-user-service/internal/service"
	"context"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpcserver.Server
}

func New(
	ctx context.Context,
	cfg config.Config,
	log *slog.Logger,
) (*App, error) {
	log.Info("connecting to postgres database", slog.String("dsn", cfg.Postgres.Dsn))
	db, err := postgrescfg.OpenDB(cfg.Postgres)
	if err != nil {
		log.Error("connecting to postgres database", logger.Err(err))

		return nil, err
	}

	jwtProvider := jwt.NewProvider(
		"secretKey",
		time.Minute*15,
		time.Hour*24,
	)

	validate := validator.New()
	err = validate.RegisterValidation("min_age", validatecfg.MinAge)
	if err != nil {
		return nil, err
	}
	err = validate.RegisterValidation("complex_password", validatecfg.ComplexPassword)
	if err != nil {
		return nil, err
	}

	userRepo := postgres.NewUserRepository(log, db)

	authService := service.NewAuthService(log, validate, jwtProvider, userRepo)

	grpcServer := grpcserver.NewServer(cfg.GRPC, log, authService)

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
