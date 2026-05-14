package handler

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	natsio "github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	protosvc "github.com/sorawaslocked/car-rental-protos/gen/service"
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"

	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
)

const depCheckTimeout = 3 * time.Second

type HealthHandler struct {
	tripsvc.UnimplementedHealthServiceServer
	log         *slog.Logger
	startTime   time.Time
	db          *sql.DB
	natsConn    *natsio.Conn
	carConn     *grpc.ClientConn
	streamConn  *grpc.ClientConn
	bookingConn *grpc.ClientConn
}

func NewHealthHandler(
	log *slog.Logger,
	db *sql.DB,
	natsConn *natsio.Conn,
	carConn *grpc.ClientConn,
	streamConn *grpc.ClientConn,
	bookingConn *grpc.ClientConn,
) *HealthHandler {
	return &HealthHandler{
		log:         pkglog.WithComponent(log, "handler.HealthHandler"),
		startTime:   time.Now(),
		db:          db,
		natsConn:    natsConn,
		carConn:     carConn,
		streamConn:  streamConn,
		bookingConn: bookingConn,
	}
}

func (h *HealthHandler) Health(ctx context.Context, _ *emptypb.Empty) (*protosvc.ServiceHealthResponse, error) {
	deps := []*protosvc.DependencyHealth{
		h.pingPostgres(ctx),
		h.pingNATS(),
		h.pingGRPCConn("car-service", h.carConn),
		h.pingGRPCConn("car-stream-service", h.streamConn),
		h.pingGRPCConn("booking-service", h.bookingConn),
	}

	status := "healthy"
	for _, d := range deps {
		if d.Status != "healthy" {
			status = "degraded"
			break
		}
	}

	return &protosvc.ServiceHealthResponse{
		Name:          "trip-service",
		Status:        status,
		Version:       "1.0.0",
		Timestamp:     timestamppb.Now(),
		UptimeSeconds: uint64(time.Since(h.startTime).Seconds()),
		Dependencies:  deps,
	}, nil
}

func (h *HealthHandler) pingPostgres(ctx context.Context) *protosvc.DependencyHealth {
	ctx, cancel := context.WithTimeout(ctx, depCheckTimeout)
	defer cancel()

	start := time.Now()
	err := h.db.PingContext(ctx)
	ms := uint32(time.Since(start).Milliseconds())

	dep := &protosvc.DependencyHealth{Name: "postgres", LatencyMs: &ms}
	if err != nil {
		dep.Status = "unhealthy"
		errStr := err.Error()
		dep.Error = &errStr
	} else {
		dep.Status = "healthy"
	}
	return dep
}

func (h *HealthHandler) pingNATS() *protosvc.DependencyHealth {
	start := time.Now()
	state := h.natsConn.Status()
	ms := uint32(time.Since(start).Milliseconds())

	dep := &protosvc.DependencyHealth{Name: "nats", LatencyMs: &ms}
	if state == natsio.CONNECTED {
		dep.Status = "healthy"
	} else {
		dep.Status = "unhealthy"
		errStr := state.String()
		dep.Error = &errStr
	}
	return dep
}

func (h *HealthHandler) pingGRPCConn(name string, conn *grpc.ClientConn) *protosvc.DependencyHealth {
	start := time.Now()
	state := conn.GetState()
	ms := uint32(time.Since(start).Milliseconds())

	dep := &protosvc.DependencyHealth{Name: name, LatencyMs: &ms}
	if state == connectivity.TransientFailure || state == connectivity.Shutdown {
		dep.Status = "unhealthy"
		errStr := state.String()
		dep.Error = &errStr
	} else {
		dep.Status = "healthy"
	}
	return dep
}
