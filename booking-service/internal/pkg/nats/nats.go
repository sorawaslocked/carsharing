package nats

import (
	"time"

	natsgo "github.com/nats-io/nats.go"
)

type Config struct {
	URL string `yaml:"url" env:"NATS_URL" env-required:"true"`
}

func NewConn(cfg Config) (*natsgo.Conn, error) {
	return natsgo.Connect(cfg.URL,
		natsgo.MaxReconnects(-1),
		natsgo.ReconnectWait(2*time.Second),
	)
}
