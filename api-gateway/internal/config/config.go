package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/grpc"
	"os"
)

type (
	Config struct {
		HTTPServer HTTPServer  `yaml:"http_server" env-required:"true"`
		GRPCServer grpc.Config `yaml:"grpc_server" env-required:"true"`
		Cookie     Cookie      `yaml:"cookie" env-required:"true"`
		Env        string      `yaml:"env" env-required:"true"`
	}

	HTTPServer struct {
		Host         string `yaml:"host" env:"HTTP_SERVER_HOST" env-required:"true"`
		Port         int    `yaml:"port" env:"HTTP_SERVER_PORT" env-required:"true"`
		ReadTimeout  string `yaml:"read_timeout" env:"HTTP_SERVER_READ_TIMEOUT" env-default:"30s"`
		WriteTimeout string `yaml:"write_timeout" env:"HTTP_SERVER_WRITE_TIMEOUT" env-default:"30s"`
		IdleTimeout  string `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
		GinMode      string `yaml:"gin_mode" env:"GIN_MODE" env-default:"debug"`
	}

	Cookie struct {
		Secure bool   `yaml:"secure" env:"COOKIE_SECURE" env-default:"false"`
		Domain string `yaml:"domain" env:"COOKIE_DOMAIN" env-default:""`
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
