package nats

import (
	pkglog "carsharing/shared/pkg/log"
	"crypto/tls"
	"log/slog"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

type PublisherConfig struct {
	URLs []string `yaml:"urls" env:"NATS_PUBLISHER_URLS" env-required:"true"`
	Name string   `yaml:"name" env:"NATS_PUBLISHER_NAME" env-required:"true"`

	Username string `yaml:"username" env:"NATS_PUBLISHER_USERNAME" env-required:"true"`
	Password string `yaml:"password" env:"NATS_PUBLISHER_PASSWORD" env-required:"true"`

	TLSServerName string `yaml:"tls_server_name" env:"NATS_PUBLISHER_TLS_SERVER_NAME"`

	Timeout             time.Duration `yaml:"timeout" env:"NATS_PUBLISHER_TIMEOUT" env-default:"5s"`
	PingInterval        time.Duration `yaml:"ping_interval" env:"NATS_PUBLISHER_PING_INTERVAL" env-default:"20s"`
	MaxPingsOutstanding int           `yaml:"max_pings_outstanding" env:"NATS_PUBLISHER_MAX_PINGS_OUTSTANDING" env-default:"3"`
}

func NewPublisher(log *slog.Logger, cfg PublisherConfig) (*nats.Conn, error) {
	log = pkglog.WithMethod(log, "nats.NewPublisher")

	opts := []nats.Option{
		nats.UserInfo(cfg.Username, cfg.Password),
		nats.Name(cfg.Name),

		nats.UserInfo(cfg.Username, cfg.Password),

		nats.Timeout(cfg.Timeout),
		nats.PingInterval(cfg.PingInterval),
		nats.MaxPingsOutstanding(cfg.MaxPingsOutstanding),

		nats.NoEcho(),
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
