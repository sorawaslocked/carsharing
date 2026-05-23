package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"carsharing/telematics-service/internal/handler"
	telematicspb "github.com/sorawaslocked/car-rental-protos/gen/service/telematics"
)

// NewGRPCServer creates a gRPC server with the telemetry streaming service registered.
// Server reflection is enabled for tooling such as grpcurl.
func NewGRPCServer(h *handler.TelematicsHandler) *grpc.Server {
	srv := grpc.NewServer()
	telematicspb.RegisterCarTelematicsStreamServiceServer(srv, h)
	reflection.Register(srv)
	return srv
}
