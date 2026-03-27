package grpc

import "time"

type Config struct {
	Host    string        `yaml:"host" env:"GRPC_SERVER_HOST" env-required:"true"`
	Port    int           `yaml:"port" env:"GRPC_SERVER_PORT" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_SERVER_TIMEOUT" env-default:"1m"`
}
