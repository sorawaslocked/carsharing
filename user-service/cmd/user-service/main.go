package main

import (
	"carsharing/user-service/internal/app"
	"carsharing/user-service/internal/config"
	pkglog "carsharing/user-service/internal/pkg/log"
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
