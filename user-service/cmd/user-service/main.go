package main

import (
	"car-rental-user-service/internal/app"
	"car-rental-user-service/internal/config"
	"car-rental-user-service/internal/pkg/logger"
	"fmt"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	application := app.New(cfg, log)
	fmt.Println(application)
}
