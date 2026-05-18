package config

import (
	"flag"
	"os"

	pkggrpc "carsharing/car-service/internal/pkg/grpc"
	pkgminio "carsharing/car-service/internal/pkg/minio"
	pkgnats "carsharing/car-service/internal/pkg/nats"
	pkgpostgres "carsharing/car-service/internal/pkg/postgres"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string                                `yaml:"env"                env:"ENV"                      env-default:"local"`
	GRPC             pkggrpc.ServerConfig                  `yaml:"grpc_server"`
	PG               pkgpostgres.Config                    `yaml:"postgres"`
	NATS             pkgnats.Config                        `yaml:"nats"`
	MinIO            pkgminio.Config                       `yaml:"minio"`
	TelematicsStream pkggrpc.TelematicsStreamServiceConfig `yaml:"telematics_stream"`
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
