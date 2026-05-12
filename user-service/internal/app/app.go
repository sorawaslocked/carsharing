package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	grpcserver "github.com/sorawaslocked/car-rental-user-service/internal/adapter/grpc"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/mailer"
	natsadapter "github.com/sorawaslocked/car-rental-user-service/internal/adapter/nats"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/postgres"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/redis"
	"github.com/sorawaslocked/car-rental-user-service/internal/config"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
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

	redisConn := rediscfg.Client(cfg.Redis)
	activationCodeCache := redis.NewActivationCodeRedisCache(redisConn)

	msMailer := mailer.New(cfg.Mailer)
	userRepo := postgres.NewUserRepository(log, db)
	publisher := natsadapter.NewPublisher(log, natsConn)

	userService := service.NewUserService(log, validate, userRepo, publisher, activationCodeCache, msMailer)

	grpcServer := grpcserver.NewServer(cfg.GRPC, log, userService)

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
