package handler

import (
	"context"
	"log/slog"
	"time"

	svcpb "github.com/sorawaslocked/car-rental-protos/gen/service"
	carsvc "github.com/sorawaslocked/car-rental-protos/gen/service/car"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HealthHandler struct {
	startTime   time.Time
	db          DBPinger
	natsChecker NATSChecker

	log *slog.Logger

	carsvc.UnimplementedHealthServiceServer
}

func NewHealthHandler(db DBPinger, natsChecker NATSChecker, log *slog.Logger) *HealthHandler {
	return &HealthHandler{
		startTime:   time.Now(),
		db:          db,
		natsChecker: natsChecker,
		log: log.With(
			slog.Group("src",
				slog.String("component", "HealthHandler"),
			),
		),
	}
}

func (h *HealthHandler) Health(ctx context.Context, _ *emptypb.Empty) (*svcpb.ServiceHealthResponse, error) {
	status := "ok"

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := h.db.Ping(pingCtx); err != nil {
		h.log.Error("postgres health check failed", slog.String("error", err.Error()))
		status = "degraded"
	}

	if !h.natsChecker.IsConnected() {
		h.log.Error("nats health check failed: not connected")
		status = "degraded"
	}

	return &svcpb.ServiceHealthResponse{
		Name:          "car-service",
		Status:        status,
		Timestamp:     timestamppb.Now(),
		UptimeSeconds: uint64(time.Since(h.startTime).Seconds()),
	}, nil
}
