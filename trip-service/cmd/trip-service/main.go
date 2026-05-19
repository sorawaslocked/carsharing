package main

import (
	pkglog "carsharing/shared/pkg/log"
	"carsharing/trip-service/internal/app"
	"carsharing/trip-service/internal/config"
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
