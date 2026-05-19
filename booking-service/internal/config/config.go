package config

import (
	"flag"
	"os"

	pkggrpc "carsharing/shared/pkg/grpc"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string                   `yaml:"env"             env:"ENV"  env-default:"local"`
	GRPC           pkggrpc.ServerConfig     `yaml:"grpc_server"`
	PG             pkgpostgres.Config       `yaml:"postgres"`
	NATSPublisher  pkgnats.PublisherConfig  `yaml:"nats_publisher"`
	NATSSubscriber pkgnats.SubscriberConfig `yaml:"nats_subscriber"`
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
