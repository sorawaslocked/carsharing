package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	grpccfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/grpc"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/mailer"
	miniocfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/minio"
	natscfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/nats"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/postgres"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/redis"
)

type Config struct {
	Env        string               `yaml:"env" env:"ENV" env-required:"true"`
	Postgres   postgres.Config      `yaml:"postgres" env-required:"true"`
	Redis      redis.Config         `yaml:"redis" env-required:"true"`
	NATS       natscfg.Config       `yaml:"nats" env-required:"true"`
	GRPC       grpccfg.ServerConfig `yaml:"grpc" env-required:"true"`
	GRPCClient grpccfg.Config       `yaml:"grpc_client"`
	Minio      miniocfg.Config      `yaml:"minio" env-required:"true"`
	Mailer     mailer.Config
}

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
