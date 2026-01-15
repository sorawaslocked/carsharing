package redis

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Host         string        `yaml:"host" env:"REDIS_HOST" required:"true"`
	Password     string        `yaml:"password" env:"REDIS_PASSWORD" required:"true"`
	User         string        `yaml:"user" env:"REDIS_USER" env-default:"car_rental"`
	MaxRetries   int           `yaml:"max_retries" env:"REDIS_MAX_RETRIES" env-default:"5"`
	DialTimeout  time.Duration `yaml:"dial_timeout" env:"REDIS_DIAL_TIMEOUT" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"REDIS_WRITE_TIMEOUT" env-default:"10s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"REDIS_READ_TIMEOUT" env-default:"10s"`
}

func Client(cfg Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Host,
		Password:     cfg.Password,
		Username:     cfg.User,
		DB:           0,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	})

	return rdb
}
