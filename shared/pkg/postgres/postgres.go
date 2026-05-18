package postgres

import (
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host string `yaml:"host"                    env:"POSTGRES_HOST"                 env-required:"true"`
	Port int    `yaml:"port"                    env:"POSTGRES_PORT"                 env-required:"true"`

	User     string `yaml:"user"                    env:"POSTGRES_USER"                 env-required:"true"`
	Password string `yaml:"password"                env:"POSTGRES_PASSWORD"             env-required:"true"`
	Database string `yaml:"database"                env:"POSTGRES_DATABASE"             env-required:"true"`

	SSLMode string `yaml:"ssl_mode"                env:"POSTGRES_SSL_MODE"             env-default:"disable"`

	MaxOpenConnections int `yaml:"max_open_connections"    env:"POSTGRES_MAX_OPEN_CONNECTIONS" env-default:"20"`
	MinOpenConnections int `yaml:"min_open_connections" env:"POSTGRES_MIN_OPEN_CONNECTIONS" env-default:"5"`

	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" env:"POSTGRES_MAX_CONN_LIFETIME" env-default:"30m"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time" env:"POSTGRES_MAX_CONN_IDLE_TIME" env-default:"10m"`

	HealthCheckPeriod time.Duration `yaml:"health_check_period" env:"POSTGRES_HEALTH_CHECK_PERIOD" env-default:"1m"`
	ConnectTimeout    time.Duration `yaml:"connect_timeout" env:"POSTGRES_CONNECT_TIMEOUT" env-default:"5s"`
}

func NewPool(log *slog.Logger, cfg Config) (*pgxpool.Pool, error) {
	log = pkglog.WithMethod(log, "postgres.NewPool")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)

	pgCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error(
			"failed to parse config",
			pkglog.Err(err),
			slog.String("host", cfg.Host),
			slog.Int("port", cfg.Port),
		)

		return nil, ErrFailedConnection
	}

	pgCfg.MaxConns = int32(cfg.MaxOpenConnections)
	pgCfg.MinConns = int32(cfg.MinOpenConnections)
	pgCfg.MaxConnLifetime = cfg.MaxConnLifetime
	pgCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	pgCfg.HealthCheckPeriod = cfg.HealthCheckPeriod
	pgCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	pool, err := pgxpool.NewWithConfig(context.Background(), pgCfg)
	if err != nil {
		log.Error(
			"failed to parse config",
			pkglog.Err(err),
			slog.String("host", cfg.Host),
			slog.Int("port", cfg.Port),
		)

		return nil, ErrFailedConnection
	}

	err = Ping(context.Background(), log, pool)
	if err != nil {
		return nil, ErrFailedConnection
	}

	return pool, nil
}

func Ping(ctx context.Context, log *slog.Logger, pool *pgxpool.Pool) error {
	log = pkglog.WithMethod(log, "postgres.Ping")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	err := pool.Ping(ctx)
	if err != nil {
		connCfg := pool.Config().ConnConfig
		log.Error(
			"pinging postgres",
			pkglog.Err(err),
			slog.String("host", connCfg.Host),
			slog.Int("port", int(connCfg.Port)),
		)

		return ErrFailedConnection
	}

	return nil
}
