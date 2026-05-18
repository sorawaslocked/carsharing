package nats

import (
	pkglog "carsharing/shared/pkg/log"
	"crypto/tls"
	"log/slog"
	"strings"

	"github.com/nats-io/nats.go"
)

type SubscriberConfig struct {
	URLs []string `yaml:"urls" env:"NATS_SUBSCRIBER_URLS" env-required:"true"`
	Name string   `yaml:"name" env:"NATS_SUBSCRIBER_NAME" env-required:"true"`

	Username string `yaml:"username" env:"NATS_SUBSCRIBER_USERNAME" env-required:"true"`
	Password string `yaml:"password" env:"NATS_SUBSCRIBER_PASSWORD" env-required:"true"`

	TLSServerName string `yaml:"tls_server_name" env:"NATS_SUBSCRIBER_TLS_SERVER_NAME"`
}

func NewSubscriber(log *slog.Logger, cfg SubscriberConfig) (*nats.Conn, error) {
	log = pkglog.WithMethod(log, "nats.NewSubscriber")

	opts := []nats.Option{
		nats.Name(cfg.Name),
		nats.UserInfo(cfg.Username, cfg.Password),
	}
	if cfg.TLSServerName != "" {
		opts = append(opts, nats.Secure(&tls.Config{
			ServerName: cfg.TLSServerName,
			MinVersion: tls.VersionTLS12,
		}))
	}
	urls := strings.Join(cfg.URLs, ",")

	nc, err := nats.Connect(urls, opts...)
	if err != nil {
		log.Error("connecting to nats", pkglog.Err(err), slog.String("urls", urls))

		return nil, ErrFailedConnection
	}

	return nc, nil
}
