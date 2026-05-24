package handler

import (
	"context"
	"log/slog"
	"time"

	servicepb "carsharing/protos/gen/service"
	servicebookingpb "carsharing/protos/gen/service/booking"
	pkglog "carsharing/shared/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	servicebookingpb.UnimplementedHealthServiceServer
	log       *slog.Logger
	version   string
	deps      map[string]Pinger
	startTime time.Time
}

func NewHealthHandler(log *slog.Logger, deps map[string]Pinger, version string) *HealthHandler {
	return &HealthHandler{
		log:       pkglog.WithComponent(log, "grpc.handler.HealthHandler"),
		version:   version,
		deps:      deps,
		startTime: time.Now(),
	}
}

func (h *HealthHandler) Health(ctx context.Context, _ *emptypb.Empty) (*servicepb.ServiceHealthResponse, error) {
	depHealths := make([]*servicepb.DependencyHealth, 0, len(h.deps))
	overallStatus := "healthy"

	for name, pinger := range h.deps {
		start := time.Now()
		err := pinger.Ping(ctx)
		latencyMs := uint32(time.Since(start).Milliseconds())

		depStatus := "healthy"
		if err != nil {
			depStatus = "unhealthy"
			overallStatus = "unhealthy"
			h.log.Error("dependency unhealthy", slog.String("dep", name), pkglog.Err(err))
		}

		depHealths = append(depHealths, &servicepb.DependencyHealth{
			Name:      name,
			Status:    depStatus,
			LatencyMs: &latencyMs,
		})
	}

	return &servicepb.ServiceHealthResponse{
		Name:          "booking-service",
		Version:       h.version,
		Status:        overallStatus,
		Timestamp:     timestamppb.Now(),
		UptimeSeconds: uint64(time.Since(h.startTime).Seconds()),
		Dependencies:  depHealths,
	}, nil
}
