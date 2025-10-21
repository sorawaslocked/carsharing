package app

import (
	"car-rental-user-service/internal/config"
	"car-rental-user-service/internal/pkg/postgres"
	"log"
)

type App struct {
}

func New(cfg config.Config) *App {
	_, err := postgres.OpenDB(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	return &App{}
}
