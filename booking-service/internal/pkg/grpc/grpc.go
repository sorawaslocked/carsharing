package grpc

import (
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerConfig struct {
	Addr string `yaml:"addr" env:"GRPC_ADDR" env-required:"true"`
}

func NewClientConn(addr string) (*googlegrpc.ClientConn, error) {
	return googlegrpc.NewClient(addr,
		googlegrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
