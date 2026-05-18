package config

import (
	"flag"
	"os"

	brevocfg "carsharing/user-service/internal/pkg/brevo"
	grpccfg "carsharing/user-service/internal/pkg/grpc"
	miniocfg "carsharing/user-service/internal/pkg/minio"
	natscfg "carsharing/user-service/internal/pkg/nats"
	"carsharing/user-service/internal/pkg/postgres"
	"carsharing/user-service/internal/pkg/redis"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string                         `yaml:"env"               env:"ENV"               env-default:"local"`
	GRPC             grpccfg.ServerConfig           `yaml:"grpc_server"`
	Postgres         postgres.Config                `yaml:"postgres"`
	Redis            redis.Config                   `yaml:"redis"`
	NATS             natscfg.Config                 `yaml:"nats"`
	Minio            miniocfg.Config                `yaml:"minio"`
	Brevo            brevocfg.Config                `yaml:"brevo"`
	DocumentAnalyzer grpccfg.DocumentAnalyzerConfig `yaml:"document_analyzer"`
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
