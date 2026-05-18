package postgres

import (
	"context"
	"log/slog"

	pkgpostgres "carsharing/shared/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Checker struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewChecker(log *slog.Logger, pool *pgxpool.Pool) *Checker {
	return &Checker{log: log, pool: pool}
}

func (c *Checker) Ping(ctx context.Context) error {
	return pkgpostgres.Ping(ctx, c.log, c.pool)
}
