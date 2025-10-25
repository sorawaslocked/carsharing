package main

import (
	"github.com/sorawaslocked/car-rental-api-gateway/internal/app"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	application := app.New(cfg, log)

	if application != nil {
		application.Run()
	}
}
