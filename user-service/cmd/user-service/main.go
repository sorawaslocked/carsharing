package main

import (
	"context"
	"github.com/sorawaslocked/car-rental-user-service/internal/app"
	"github.com/sorawaslocked/car-rental-user-service/internal/config"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/logger"
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
