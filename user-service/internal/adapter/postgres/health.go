package postgres

import (
	"context"
	"database/sql"
)

type Checker struct {
	db *sql.DB
}

func NewChecker(db *sql.DB) *Checker {
	return &Checker{db: db}
}

func (c *Checker) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}
