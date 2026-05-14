package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sorawaslocked/car-rental-booking-service/internal/app"
	"github.com/sorawaslocked/car-rental-booking-service/internal/config"
	pkglog "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/log"
)

func main() {
	cfg := config.MustLoad()
	log := pkglog.SetupLogger(cfg.Env)

	log.Info("starting booking service", slog.String("env", cfg.Env))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := app.New(ctx, log, cfg)
	if err != nil {
		log.Error("failed to initialise application", pkglog.Err(err))
		os.Exit(1)
	}

	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			log.Error("grpc server stopped", pkglog.Err(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("shutting down")
	cancel()
	application.GRPCServer.Stop()
	log.Info("booking service stopped")
}
