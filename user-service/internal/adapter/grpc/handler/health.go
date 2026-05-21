package handler

import (
	"context"
	"log/slog"
	"time"

	pkglog "carsharing/shared/pkg/log"

	servicepb "github.com/sorawaslocked/car-rental-protos/gen/service"
	usersvc "github.com/sorawaslocked/car-rental-protos/gen/service/user"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	log       *slog.Logger
	deps      map[string]Pinger
	startTime time.Time
	usersvc.UnimplementedHealthServiceServer
}

func NewHealthHandler(log *slog.Logger, deps map[string]Pinger) *HealthHandler {
	return &HealthHandler{
		log:       pkglog.WithComponent(log, "grpc.HealthHandler"),
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

		status := "healthy"
		if err != nil {
			status = "unhealthy"
			overallStatus = "unhealthy"
			h.log.Error("dependency unhealthy", slog.String("dep", name), pkglog.Err(err))
		}

		depHealths = append(depHealths, &servicepb.DependencyHealth{
			Name:      name,
			Status:    status,
			LatencyMs: &latencyMs,
		})
	}

	return &servicepb.ServiceHealthResponse{
		Name:          "user-service",
		Status:        overallStatus,
		Timestamp:     timestamppb.Now(),
		UptimeSeconds: uint64(time.Since(h.startTime).Seconds()),
		Dependencies:  depHealths,
	}, nil
}
