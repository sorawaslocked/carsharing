package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	grpcserver "github.com/sorawaslocked/car-rental-booking-service/internal/adapter/grpc"
	pkgnats "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/nats"
	pkgpostgres "github.com/sorawaslocked/car-rental-booking-service/internal/pkg/postgres"
)

type Config struct {
	Env        string             `yaml:"env"         env:"ENV" env-required:"true"`
	GRPCServer grpcserver.Config  `yaml:"grpc_server"`
	Postgres   pkgpostgres.Config `yaml:"postgres"`
	NATS       pkgnats.Config     `yaml:"nats"`
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
