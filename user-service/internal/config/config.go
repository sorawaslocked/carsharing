package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/grpc"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/postgres"
	"os"
)

type (
	Config struct {
		Env      string          `yaml:"env" env:"ENV" required:"true"`
		Postgres postgres.Config `yaml:"postgres" env-required:"true"`
		GRPC     grpc.Config     `yaml:"grpc" env-required:"true"`
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
