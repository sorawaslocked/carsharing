package grpc

import (
	"context"
	"log/slog"

	pkggrpc "carsharing/shared/pkg/grpc"
	"google.golang.org/grpc"
)

type Pinger struct {
	log  *slog.Logger
	conn *grpc.ClientConn
}

func NewPinger(log *slog.Logger, conn *grpc.ClientConn) *Pinger {
	return &Pinger{log: log, conn: conn}
}

func (p *Pinger) Ping(ctx context.Context) error {
	_, err := pkggrpc.PingClient(ctx, p.log, p.conn)
	return err
}
