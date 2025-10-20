package app

import (
	"context"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc"
	httpserver "github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	grpcconn "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/grpc"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/service"
	grpcSvc "github.com/sorawaslocked/car-rental-protos/gen/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	httpServer *httpserver.Server
}

func New(cfg config.Config) *App {
	log.Println("starting api-gateway")

	log.Println("connecting to grpc server")
	authServiceGrpcConn, err := grpcconn.Connect(
		cfg.GRPCServer.Client.AuthServiceURL,
		cfg.GRPCServer.Client,
	)
	if err != nil {
		log.Println("could not connect to auth service", err.Error())
	}

	err = grpcconn.PingServer(authServiceGrpcConn)
	if err != nil {
		log.Println("could not connect to auth service", err.Error())
	}

	authServiceGrpcClient := grpcSvc.NewAuthServiceClient(authServiceGrpcConn)
	authServiceGrpcHandler := grpc.NewAuthHandler(authServiceGrpcClient)
	authService := service.NewAuthService(authServiceGrpcHandler)

	httpServer := httpserver.New(cfg.HTTPServer, authService)

	app := &App{
		httpServer: httpServer,
	}

	return app
}

func (a *App) Run() {
	a.httpServer.MustRun()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	s := <-shutdown
	log.Println("received shutdown signal", s)
	a.Stop()
	log.Println("graceful shutdown complete")
}

func (a *App) Stop() {
	log.Println("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.httpServer.Stop(ctx)
}
