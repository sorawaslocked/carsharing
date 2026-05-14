package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sorawaslocked/car-rental-trip-service/internal/adapter/grpc/interceptor"
)

func NewClientConn(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(interceptor.MetadataForwardingUnaryInterceptor),
		grpc.WithChainStreamInterceptor(interceptor.MetadataForwardingStreamInterceptor),
	)
}
