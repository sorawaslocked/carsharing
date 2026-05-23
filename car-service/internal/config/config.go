package config

import (
	"flag"
	"os"
	"time"

	pkggrpc "carsharing/shared/pkg/grpc"
	pkgminio "carsharing/shared/pkg/minio"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type TelemetryConfig struct {
	StalenessThreshold time.Duration `yaml:"staleness_threshold" env-default:"2m"`
}

type Config struct {
	Env             string                   `yaml:"env"              env:"ENV"             env-default:"local"`
	GRPC            pkggrpc.ServerConfig     `yaml:"grpc_server"`
	PG              pkgpostgres.Config       `yaml:"postgres"`
	NATSPublisher   pkgnats.PublisherConfig  `yaml:"nats_publisher"`
	NATSSubscriber  pkgnats.SubscriberConfig `yaml:"nats_subscriber"`
	MinIO           pkgminio.Config          `yaml:"minio"`
	TelemetryStream pkggrpc.ClientConfig     `yaml:"telemetry_stream"`
	Telemetry       TelemetryConfig          `yaml:"telemetry"`
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
