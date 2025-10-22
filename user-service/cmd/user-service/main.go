package main

import (
	"car-rental-user-service/internal/app"
	"car-rental-user-service/internal/config"
	"car-rental-user-service/internal/pkg/logger"
	"context"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	application, err := app.New(ctx, cfg, log)
	if err != nil {
		panic(err)
	}
	application.Run()
}
