package handler

import (
	"context"
	"database/sql"
	"log/slog"

	pkglog "carsharing/booking-service/internal/pkg/log"
	natsgo "github.com/nats-io/nats.go"
	servicepb "github.com/sorawaslocked/car-rental-protos/gen/service"
	servicebookingpb "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HealthHandler struct {
	servicebookingpb.UnimplementedHealthServiceServer
	log *slog.Logger
	db  *sql.DB
	nc  *natsgo.Conn
}

func NewHealthHandler(log *slog.Logger, db *sql.DB, nc *natsgo.Conn) *HealthHandler {
	return &HealthHandler{
		log: pkglog.WithComponent(log, "grpc.HealthHandler"),
		db:  db,
		nc:  nc,
	}
}

func (h *HealthHandler) Health(ctx context.Context, _ *emptypb.Empty) (*servicepb.ServiceHealthResponse, error) {
	log := pkglog.WithMethod(h.log, "Health")

	deps := make([]*servicepb.DependencyHealth, 0, 2)

	pgStatus := "healthy"
	if err := h.db.PingContext(ctx); err != nil {
		log.Error("postgres health check failed", pkglog.Err(err))
		pgStatus = "unhealthy"
	}
	deps = append(deps, &servicepb.DependencyHealth{Name: "postgres", Status: pgStatus})

	natsStatus := "healthy"
	if h.nc.Status() != natsgo.CONNECTED {
		natsStatus = "unhealthy"
		log.Error("nats health check failed", slog.String("status", h.nc.Status().String()))
	}
	deps = append(deps, &servicepb.DependencyHealth{Name: "nats", Status: natsStatus})

	return &servicepb.ServiceHealthResponse{
		Name:         "booking-service",
		Status:       "healthy",
		Timestamp:    timestamppb.Now(),
		Dependencies: deps,
	}, nil
}
