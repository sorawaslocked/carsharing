package logger

import (
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/logger/pretty"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettyLogger()
	case envDev:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		))
	case envProd:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		))
	}

	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := pretty.HandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.PrettyHandler(os.Stdout)

	return slog.New(handler)
}

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}
