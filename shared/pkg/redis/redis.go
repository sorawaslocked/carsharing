package redis

import (
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
	Port int    `yaml:"port" env:"REDIS_PORT" env-required:"true"`

	Username string `yaml:"username" env:"REDIS_USERNAME" env-required:"true"`
	Password string `yaml:"password"  env:"REDIS_PASSWORD" env-required:"true"`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`

	TLSServerName string `yaml:"tls_server_name" env:"REDIS_TLS_SERVER_NAME"`

	PoolSize        int           `yaml:"pool_size" env:"REDIS_POOL_SIZE" env-default:"10"`
	MinIdleConns    int           `yaml:"min_idle_conns" env:"REDIS_MIN_IDLE_CONNS" env-default:"3"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"REDIS_MAX_IDLE_CONNS" env-default:"10"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env:"REDIS_CONN_MAX_IDLE_TIME" env-default:"5m"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"REDIS_CONN_MAX_LIFETIME" env-default:"30m"`

	DialTimeout  time.Duration `yaml:"dial_timeout" env:"REDIS_DIAL_TIMEOUT" env-default:"5s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"REDIS_READ_TIMEOUT" env-default:"3s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"REDIS_WRITE_TIMEOUT" env-default:"3s"`
	PoolTimeout  time.Duration `yaml:"pool_timeout" env:"REDIS_POOL_TIMEOUT" env-default:"4s"`
}

func NewClient(log *slog.Logger, cfg Config) (*redis.Client, error) {
	log = pkglog.WithMethod(log, "redis.NewClient")

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	opts := &redis.Options{
		Addr:     addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,

		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,

		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
	}
	if cfg.TLSServerName != "" {
		opts.TLSConfig = &tls.Config{
			ServerName: cfg.TLSServerName,
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(opts)

	err := PingClient(context.Background(), log, client)
	if err != nil {
		return nil, ErrFailedConnection
	}

	return client, nil
}

func PingClient(ctx context.Context, log *slog.Logger, client *redis.Client) error {
	log = pkglog.WithMethod(log, "redis.PingClient")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Error("pinging redis", pkglog.Err(err))

		return ErrFailedConnection
	}

	return nil
}
