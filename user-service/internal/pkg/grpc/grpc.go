package grpc

type ServerConfig struct {
	Addr string `yaml:"addr" env:"GRPC_ADDR" env-required:"true"`
}
