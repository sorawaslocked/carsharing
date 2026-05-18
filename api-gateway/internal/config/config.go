package config

import (
	"flag"
	"os"
	"time"

	pkggrpc "carsharing/shared/pkg/grpc"
	pkgjwt "carsharing/shared/pkg/jwt"
	pkgnats "carsharing/shared/pkg/nats"
	pkgredis "carsharing/shared/pkg/redis"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTPServer HTTPServer               `yaml:"http_server" env-required:"true"`
		GRPCServer GRPCServer               `yaml:"grpc_server" env-required:"true"`
		Redis      pkgredis.Config          `yaml:"redis" env-required:"true"`
		Cache      CacheConfig              `yaml:"cache"`
		NATS       pkgnats.SubscriberConfig `yaml:"nats" env-required:"true"`
		JWT        pkgjwt.Config            `yaml:"jwt" env-required:"true"`
		Env        string                   `yaml:"env" env-required:"true"`
	}

	GRPCServer struct {
		UserService    pkggrpc.ClientConfig `yaml:"user_service" env-required:"true"`
		CarService     pkggrpc.ClientConfig `yaml:"car_service" env-required:"true"`
		BookingService pkggrpc.ClientConfig `yaml:"booking_service" env-required:"true"`
		TripService    pkggrpc.ClientConfig `yaml:"trip_service" env-required:"true"`
	}

	CacheConfig struct {
		MetadataTTL time.Duration `yaml:"metadata_ttl" env:"CACHE_METADATA_TTL" env-default:"1h"`
		SessionTTL  time.Duration `yaml:"session_ttl" env:"CACHE_SESSION_TTL" env-default:"24h"`
	}

	HTTPServer struct {
		Host         string `yaml:"host" env:"HTTP_SERVER_HOST" env-required:"true"`
		Port         int    `yaml:"port" env:"HTTP_SERVER_PORT" env-required:"true"`
		ReadTimeout  string `yaml:"read_timeout" env:"HTTP_SERVER_READ_TIMEOUT" env-default:"30s"`
		WriteTimeout string `yaml:"write_timeout" env:"HTTP_SERVER_WRITE_TIMEOUT" env-default:"30s"`
		IdleTimeout  string `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
		GinMode      string `yaml:"gin_mode" env:"GIN_MODE" env-default:"debug"`
		Cookie       Cookie `yaml:"cookie" env-required:"true"`
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
