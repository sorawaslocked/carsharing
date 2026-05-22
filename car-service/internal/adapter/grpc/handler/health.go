package handler

import (
	"context"
	"log/slog"
	"time"

	pkglog "carsharing/shared/pkg/log"
	svcpb "github.com/sorawaslocked/car-rental-protos/gen/service"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HealthHandler struct {
	log       *slog.Logger
	deps      map[string]Pinger
	startTime time.Time
	carsvc.UnimplementedHealthServiceServer
}

func NewHealthHandler(log *slog.Logger, deps map[string]Pinger) *HealthHandler {
	return &HealthHandler{
		log:       pkglog.WithComponent(log, "grpc.handler.HealthHandler"),
		deps:      deps,
		startTime: time.Now(),
	}
}

func (h *HealthHandler) Health(ctx context.Context, _ *emptypb.Empty) (*svcpb.ServiceHealthResponse, error) {
	overallStatus := "healthy"

	for name, pinger := range h.deps {
		if err := pinger.Ping(ctx); err != nil {
			overallStatus = "unhealthy"
			h.log.Error("dependency unhealthy", slog.String("dep", name), pkglog.Err(err))
		}
	}

	return &svcpb.ServiceHealthResponse{
		Name:          "car-service",
		Status:        overallStatus,
		Timestamp:     timestamppb.Now(),
		UptimeSeconds: uint64(time.Since(h.startTime).Seconds()),
	}, nil
}
