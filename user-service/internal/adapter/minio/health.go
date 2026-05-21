package minio

import (
	pkglog "carsharing/shared/pkg/log"
	pkgminio "carsharing/shared/pkg/minio"
	"context"
	"log/slog"

	"github.com/minio/minio-go/v7"
)

type Pinger struct {
	log    *slog.Logger
	client *minio.Client
}

func NewPinger(log *slog.Logger, client *minio.Client) *Pinger {
	return &Pinger{
		client: client,
		log:    pkglog.WithComponent(log, "minio.Pinger"),
	}
}

func (c *Pinger) Ping(ctx context.Context) error {
	return pkgminio.Ping(ctx, c.log, c.client)
}
