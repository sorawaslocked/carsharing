package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/handler"
	httpserver "github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	grpcconn "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/grpc"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/logger"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/service"
	authsvc "github.com/sorawaslocked/car-rental-protos/gen/service/auth"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
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
	authServiceGrpcHandler := handler.NewAuthHandler(authServiceGrpcClient)
	authService := service.NewAuthService(authServiceGrpcHandler)

	userServiceGrpcClient := usersvc.NewUserServiceClient(userServiceGrpcConn)
	userServiceGrpcHandler := handler.NewUserHandler(userServiceGrpcClient)
	userService := service.NewUserService(userServiceGrpcHandler)

	carModelServiceGrpcHandler := handler.NewCarModelHandler()
	carModelService := service.NewCarModelService(carModelServiceGrpcHandler)

	carServiceGrpcHandler := handler.NewCarHandler()
	carService := service.NewCarService(carServiceGrpcHandler)

	carInsuranceServiceGrpcHandler := handler.NewInsuranceHandler()
	carInsuranceService := service.NewCarInsuranceService(carInsuranceServiceGrpcHandler)

	carMaintenanceServiceGrpcHandler := handler.NewCarMaintenanceHandler()
	carMaintenanceService := service.NewCarMaintenanceService(carMaintenanceServiceGrpcHandler)

	zoneServiceGrpcHandler := handler.NewZoneHandler()
	zoneService := service.NewZoneService(zoneServiceGrpcHandler)

	httpServer := httpserver.New(
		cfg.HTTPServer,
		cfg.Cookie,
		log,
		authService,
		userService,
		carModelService,
		carService,
		carInsuranceService,
		carMaintenanceService,
		zoneService,
	)

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
