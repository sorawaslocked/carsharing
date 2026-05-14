package handler

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	protosvc "github.com/sorawaslocked/car-rental-protos/gen/service"
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"

	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
)

type HealthHandler struct {
	tripsvc.UnimplementedHealthServiceServer
	log       *slog.Logger
	startTime time.Time
}

func NewHealthHandler(log *slog.Logger) *HealthHandler {
	return &HealthHandler{
		log:       pkglog.WithComponent(log, "handler.HealthHandler"),
		startTime: time.Now(),
	}
}

func (h *HealthHandler) Health(_ context.Context, _ *emptypb.Empty) (*protosvc.ServiceHealthResponse, error) {
	return &protosvc.ServiceHealthResponse{
		Name:          "trip-service",
		Status:        "healthy",
		Version:       "1.0.0",
		Timestamp:     timestamppb.Now(),
		UptimeSeconds: uint64(time.Since(h.startTime).Seconds()),
	}, nil
}
