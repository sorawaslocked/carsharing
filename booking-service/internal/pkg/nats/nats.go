package nats

import (
	"fmt"

	natsgo "github.com/nats-io/nats.go"
)

type Config struct {
	URL string `yaml:"url" env:"NATS_URL" env-required:"true"`
}

func New(cfg Config) (*natsgo.Conn, error) {
	nc, err := natsgo.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	return nc, nil
}
