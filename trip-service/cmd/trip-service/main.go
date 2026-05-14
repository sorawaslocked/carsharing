package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sorawaslocked/car-rental-trip-service/internal/app"
	"github.com/sorawaslocked/car-rental-trip-service/internal/config"
	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
)

func main() {
	cfg := config.MustLoad()
	log := pkglog.SetupLogger(cfg.Env)

	a, err := app.New(log, cfg)
	if err != nil {
		log.Error("failed to initialize app", pkglog.Err(err))
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.Run(); err != nil {
			log.Error("server stopped with error", pkglog.Err(err))
			quit <- syscall.SIGTERM
		}
	}()

	sig := <-quit
	log.Info("received signal, shutting down", slog.String("signal", sig.String()))
	a.Stop()
}
