package app

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc"
	httpserver "github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	grpcconn "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/grpc"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/logger"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/service"
	authsvc "github.com/sorawaslocked/car-rental-protos/gen/service/auth"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg        config.Config
	log        *slog.Logger
	httpServer *httpserver.Server
}

func New(cfg config.Config, log *slog.Logger) *App {
	log = log.With(slog.String("appId", "api-gateway"))

	userServiceLog := log.With(
		slog.String("grpcServiceId", "user-service"),
		slog.String("grpcURL", cfg.GRPCServer.Client.UserServiceURL),
	)

	userServiceLog.Info("connecting to grpc server")
	userServiceLog.Debug(
		"connecting to grpc server",
		slog.Int("maxReceiveSizeMb", cfg.GRPCServer.Client.MaxReceiveSizeMb),
		slog.Duration("timeout", cfg.GRPCServer.Client.Timeout),
		slog.Duration("timeKeepAlive", cfg.GRPCServer.Client.TimeKeepAlive),
	)

	userServiceGrpcConn, err := grpcconn.Connect(
		cfg.GRPCServer.Client.UserServiceURL,
		cfg.GRPCServer.Client,
	)
	if err != nil {
		userServiceLog.Error("error connecting to grpc server")
		userServiceLog.Debug("error connecting to grpc server", logger.Err(err))

		return nil
	}

	err = grpcconn.PingServer(userServiceGrpcConn)
	if err != nil {
		userServiceLog.Error("error pinging grpc server")
		userServiceLog.Debug("error pinging grpc server", logger.Err(err))
	}

	authServiceGrpcClient := authsvc.NewAuthServiceClient(userServiceGrpcConn)
	authServiceGrpcHandler := grpc.NewAuthHandler(authServiceGrpcClient)
	authService := service.NewAuthService(authServiceGrpcHandler)

	httpServer := httpserver.New(cfg.Env, cfg.HTTPServer, log, authService)

	app := &App{
		cfg:        cfg,
		log:        log,
		httpServer: httpServer,
	}

	return app
}

func (a *App) Run() {
	a.httpServer.MustRun()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	s := <-shutdown

	a.log.Info("received shutdown signal", slog.String("signal", s.String()))

	a.Stop()

	a.log.Info("graceful shutdown complete")
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.httpServer.Stop(ctx)
}
