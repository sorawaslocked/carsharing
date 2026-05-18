package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerConfig struct {
	Addr string `yaml:"addr" env:"GRPC_ADDR" env-required:"true"`
}

type TelematicsStreamServiceConfig struct {
	Addr string `yaml:"addr" env:"TELEMATICS_STREAM_ADDR" env-required:"true"`
}

func NewClientConn(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
