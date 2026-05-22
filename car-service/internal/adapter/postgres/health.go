package postgres

import (
	"context"
	"log/slog"

	pkgpostgres "carsharing/shared/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pinger struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewPinger(log *slog.Logger, pool *pgxpool.Pool) *Pinger {
	return &Pinger{log: log, pool: pool}
}

func (p *Pinger) Ping(ctx context.Context) error {
	return pkgpostgres.Ping(ctx, p.log, p.pool)
}
