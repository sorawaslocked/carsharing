package redis

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
)

type Config struct {
	Address      string        `yaml:"address" env:"REDIS_ADDRESS" env-required:"true"`
	Password     string        `yaml:"password"  env:"REDIS_PASSWORD" env-required:"true"`
	DB           int           `yaml:"db" env:"REDIS_DB" env-default:"0"`
	PoolSize     int           `yaml:"pool_size" env:"REDIS_POOL_SIZE" env-default:"10"`
	DialTimeout  time.Duration `yaml:"dial_timeout" env:"REDIS_DIAL_TIMEOUT" env-default:"5s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"REDIS_READ_TIMEOUT" env-default:"3s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"REDIS_WRITE_TIMEOUT" env-default:"3s"`
	MetadataTTL  time.Duration `yaml:"metadata_ttl" env:"REDIS_METADATA_TTL" env-default:"1h"`
	SessionTTL   time.Duration `yaml:"session_ttl" env:"REDIS_SESSION_TTL" env-default:"24h"`
}

func NewClient(ctx context.Context, config *Config, oldLog *slog.Logger) (*redis.Client, error) {
	logger := oldLog.With(
		slog.Group("src",
			slog.String("method", "redis.NewClient"),
		),
		slog.Group("redis",
			slog.String("address", config.Address),
			slog.Int("db", config.DB),
		),
	)

	opts := &redis.Options{
		Addr:         config.Address,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	logger.Info("connecting to redis")
	client := redis.NewClient(opts)

	logger.Info("pinging redis")
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Error("pinging redis", log.Err(err))

		return nil, ErrFailedConnection
	}

	return client, nil
}
