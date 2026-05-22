package minio

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"
	pkgminio "carsharing/shared/pkg/minio"
	"github.com/minio/minio-go/v7"
)

type Pinger struct {
	log    *slog.Logger
	client *minio.Client
}

func NewPinger(log *slog.Logger, client *minio.Client) *Pinger {
	return &Pinger{
		log:    pkglog.WithComponent(log, "minio.Pinger"),
		client: client,
	}
}

func (p *Pinger) Ping(ctx context.Context) error {
	return pkgminio.Ping(ctx, p.log, p.client)
}
