package config

import (
	"flag"
	"os"

	pkggrpc "carsharing/shared/pkg/grpc"
	pkgminio "carsharing/shared/pkg/minio"
	pkgnats "carsharing/shared/pkg/nats"
	pkgpostgres "carsharing/shared/pkg/postgres"
	pkgredis "carsharing/shared/pkg/redis"
	brevocfg "carsharing/user-service/internal/pkg/brevo"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string                  `yaml:"env"               env:"ENV"               env-default:"local"`
	Version          string                  `yaml:"version"           env:"VERSION"           env-default:"1.0.0"`
	GRPC             pkggrpc.ServerConfig    `yaml:"grpc_server"`
	Postgres         pkgpostgres.Config      `yaml:"postgres"`
	Redis            pkgredis.Config         `yaml:"redis"`
	NATS             pkgnats.PublisherConfig `yaml:"nats"`
	Minio            pkgminio.Config         `yaml:"minio"`
	Brevo            brevocfg.Config         `yaml:"brevo"`
	DocumentAnalyzer pkggrpc.ClientConfig    `yaml:"document_analyzer"`
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
