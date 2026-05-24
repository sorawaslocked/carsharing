package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	telemetrypb "carsharing/protos/gen/service/telemetry"
	"carsharing/telematics-service/internal/handler"
)

// NewGRPCServer creates a gRPC server with the telemetry streaming service registered.
// Server reflection is enabled for tooling such as grpcurl.
func NewGRPCServer(h *handler.TelematicsHandler) *grpc.Server {
	srv := grpc.NewServer()
	telemetrypb.RegisterCarTelemetryStreamServiceServer(srv, h)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, healthSrv)
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(srv)
	return srv
}
