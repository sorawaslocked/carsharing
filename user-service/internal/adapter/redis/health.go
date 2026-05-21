package redis

import (
	"context"
	"log/slog"

	pkgredis "carsharing/shared/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type Pinger struct {
	log    *slog.Logger
	client *redis.Client
}

func NewPinger(log *slog.Logger, client *redis.Client) *Pinger {
	return &Pinger{log: log, client: client}
}

func (p *Pinger) Ping(ctx context.Context) error {
	return pkgredis.PingClient(ctx, p.log, p.client)
}
