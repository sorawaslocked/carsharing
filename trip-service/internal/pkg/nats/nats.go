package nats

import (
	"time"

	natsio "github.com/nats-io/nats.go"
)

type Config struct {
	URL string `yaml:"url" env:"NATS_URL" env-required:"true"`
}

func NewConn(cfg Config) (*natsio.Conn, error) {
	return natsio.Connect(cfg.URL,
		natsio.MaxReconnects(-1),
		natsio.ReconnectWait(2*time.Second),
	)
}
