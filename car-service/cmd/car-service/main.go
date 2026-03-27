package main

import (
	"github.com/sorawaslocked/car-rental-car-service/internal/app"
	"github.com/sorawaslocked/car-rental-car-service/internal/config"
	"github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
)

func main() {
	cfg := config.MustLoad()

	logger := log.SetupLogger(cfg.Env)

	application, err := app.New(cfg, logger)
	if err != nil {
		panic(err)
	}
	application.Run()
}
