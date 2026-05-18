package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	telematicspb "github.com/sorawaslocked/car-rental-protos/gen/service/telematics"
	"github.com/sorawaslocked/car-rental-telematics/internal/handler"
)

// NewGRPCServer creates a gRPC server with the telematics streaming service registered.
// Server reflection is enabled for tooling such as grpcurl.
func NewGRPCServer(h *handler.TelematicsHandler) *grpc.Server {
	srv := grpc.NewServer()
	telematicspb.RegisterCarTelematicsStreamServiceServer(srv, h)
	reflection.Register(srv)
	return srv
}
