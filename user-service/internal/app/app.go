package app

import (
	"car-rental-user-service/internal/config"
	"car-rental-user-service/internal/pkg/logger"
	"car-rental-user-service/internal/pkg/postgres"
	"log/slog"
)

type App struct {
	log *slog.Logger
}

func New(cfg config.Config, log *slog.Logger) *App {
	_, err := postgres.OpenDB(cfg.Postgres)
	if err != nil {
		log.Error("opening DB", logger.Err(err))

		return nil
	}

	return &App{}
}
