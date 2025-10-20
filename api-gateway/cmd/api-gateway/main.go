package main

import (
	"github.com/sorawaslocked/car-rental-api-gateway/internal/app"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application := app.New(cfg)

	application.Run()
}
