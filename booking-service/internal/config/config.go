package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env"         env:"ENV"          env-required:"true"`
	GRPCServer GRPCServer `yaml:"grpc_server"`
	Postgres   Postgres   `yaml:"postgres"`
	NATS       NATS       `yaml:"nats"`
}

type GRPCServer struct {
	Address string `yaml:"address" env:"GRPC_ADDRESS" env-default:"0.0.0.0:9997"`
}

type Postgres struct {
	Host     string `yaml:"host"     env:"PG_HOST"     env-required:"true"`
	Port     int    `yaml:"port"     env:"PG_PORT"     env-default:"5432"`
	User     string `yaml:"user"     env:"PG_USER"     env-required:"true"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
	Database string `yaml:"database" env:"PG_DATABASE" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env:"PG_SSL_MODE" env-default:"disable"`
}

type NATS struct {
	URL string `yaml:"url" env:"NATS_URL" env-required:"true"`
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
