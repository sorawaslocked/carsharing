package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
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

func Client(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	})
}
