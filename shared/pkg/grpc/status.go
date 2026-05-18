package grpc

import "google.golang.org/grpc/health/grpc_health_v1"

const (
	ServerStatusServing    = "serving"
	ServerStatusNotServing = "not_serving"
	ServerStatusUnknown    = "unknown"
)

func convertHealthStatus(hs grpc_health_v1.HealthCheckResponse_ServingStatus) string {
	switch hs {
	case grpc_health_v1.HealthCheckResponse_SERVING:
		return ServerStatusServing
	case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
		return ServerStatusNotServing
	case grpc_health_v1.HealthCheckResponse_UNKNOWN:
		return ServerStatusUnknown
	default:
		return ServerStatusUnknown
	}
}
