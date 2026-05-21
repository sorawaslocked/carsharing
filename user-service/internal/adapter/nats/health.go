package nats

import (
	"context"
	"log/slog"

	pkgnats "carsharing/shared/pkg/nats"
	"github.com/nats-io/nats.go"
)

type Pinger struct {
	log  *slog.Logger
	conn *nats.Conn
}

func NewPinger(log *slog.Logger, conn *nats.Conn) *Pinger {
	return &Pinger{log: log, conn: conn}
}

func (p *Pinger) Ping(ctx context.Context) error {
	return pkgnats.Ping(ctx, p.log, p.conn)
}
