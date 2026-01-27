package postgres

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type Config struct {
	Dsn                string        `yaml:"dsn" env:"POSTGRES_DSN" env-required:"true"`
	MaxOpenConnections int           `yaml:"max_open_connections" env:"POSTGRES_MAX_OPEN_CONNECTIONS" env-default:"25"`
	MaxIdleConnections int           `yaml:"max_idle_connections" env:"POSTGRES_MAX_IDLE_CONNECTIONS" env-default:"25"`
	MaxIdleTime        time.Duration `yaml:"max_idle_time" env:"POSTGRES_MAX_IDLE_TIME" env-default:"15m"`
}

func OpenDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
