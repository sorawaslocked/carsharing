package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"carsharing/trip-service/internal/adapter/grpc/interceptor"
)

type ServerConfig struct {
	Addr string `yaml:"addr" env:"GRPC_ADDR" env-required:"true"`
}

type CarServiceConfig struct {
	Addr string `yaml:"addr" env:"CAR_SERVICE_ADDR" env-required:"true"`
}

type CarStreamServiceConfig struct {
	Addr string `yaml:"addr" env:"CAR_STREAM_SERVICE_ADDR" env-required:"true"`
}

type BookingServiceConfig struct {
	Addr string `yaml:"addr" env:"BOOKING_SERVICE_ADDR" env-required:"true"`
}

func NewClientConn(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptor.MetadataForwardingUnaryInterceptor),
		grpc.WithChainStreamInterceptor(interceptor.MetadataForwardingStreamInterceptor),
	)
}
