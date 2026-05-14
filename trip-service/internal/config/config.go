package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string     `yaml:"env"          env:"ENV"          env-default:"local"`
	GRPC GRPCConfig `yaml:"grpc_server"`
	PG   PGConfig   `yaml:"postgres"`
	NATS NATSConfig `yaml:"nats"`

	CarService       ClientConfig `yaml:"car_service"`
	CarStreamService ClientConfig `yaml:"car_stream_service"`
	BookingService   ClientConfig `yaml:"booking_service"`
}

type GRPCConfig struct {
	Port int `yaml:"port" env:"GRPC_PORT" env-default:"9996"`
}

type PGConfig struct {
	DSN string `yaml:"dsn" env:"PG_DSN" env-required:"true"`
}

type NATSConfig struct {
	URL string `yaml:"url" env:"NATS_URL" env-default:"nats://localhost:4222"`
}

type ClientConfig struct {
	Addr string `yaml:"addr"`
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
