package nats

import (
	"log/slog"

	pkglog "carsharing/api-gateway/internal/pkg/log"
	nc "github.com/nats-io/nats.go"
)

type Config struct {
	URL string `yaml:"url" env:"NATS_URL" env-required:"true"`
}

func Connect(cfg Config, logger *slog.Logger) (*nc.Conn, error) {
	log := logger.With(
		slog.Group("src", slog.String("method", "nats.Connect")),
		slog.String("natsURL", cfg.URL),
	)

	log.Info("connecting to nats")

	conn, err := nc.Connect(cfg.URL)
	if err != nil {
		log.Error("connecting to nats", pkglog.Err(err))

		return nil, ErrConnectionFailed
	}

	return conn, nil
}
