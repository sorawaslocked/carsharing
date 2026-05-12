package nats

import (
	"github.com/nats-io/nats.go"
)

type Config struct {
	URL string `yaml:"url" env:"NATS_URL" env-required:"true"`
}

func Connect(cfg Config) (*nats.Conn, error) {
	return nats.Connect(cfg.URL)
}
