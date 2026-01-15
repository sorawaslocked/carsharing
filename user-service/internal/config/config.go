package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/grpc"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/jwt"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/postgres"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/redis"
	"os"
)

type (
	Config struct {
		Env              string          `yaml:"env" env:"ENV" env-required:"true"`
		Postgres         postgres.Config `yaml:"postgres" env-required:"true"`
		Redis            redis.Config    `yaml:"redis" env-required:"true"`
		GRPC             grpc.Config     `yaml:"grpc" env-required:"true"`
		JWT              jwt.Config      `yaml:"jwt" env-required:"true"`
		MailerSendAPIKey string          `env:"MAILER_SEND_API_KEY" env-required:"true"`
	}
)

func MustLoad() Config {
	cfgPath := fetchConfigPath()

	if cfgPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		panic("config file does not exist: " + cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "config file path")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
