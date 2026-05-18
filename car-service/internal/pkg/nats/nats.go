package nats

import (
	"time"

	"github.com/nats-io/nats.go"
)

type Config struct {
	URL string `yaml:"url" env:"NATS_URL" env-required:"true"`
}

func NewConn(cfg Config) (*nats.Conn, error) {
	return nats.Connect(cfg.URL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
}
