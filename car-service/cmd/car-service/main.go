package main

import (
	"github.com/sorawaslocked/car-rental-car-service/internal/app"
	"github.com/sorawaslocked/car-rental-car-service/internal/config"
	pkglog "github.com/sorawaslocked/car-rental-car-service/internal/pkg/log"
)

func main() {
	cfg := config.MustLoad()
	log := pkglog.SetupLogger(cfg.Env)

	a, err := app.New(log, cfg)
	if err != nil {
		log.Error("failed to initialize app", pkglog.Err(err))
		panic(err)
	}

	if err = a.Run(); err != nil {
		log.Error("app run failed", pkglog.Err(err))
		panic(err)
	}
}
