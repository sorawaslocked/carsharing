package handler

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	basetripmb "github.com/sorawaslocked/car-rental-protos/gen/base/trip"
	tripsvc "github.com/sorawaslocked/car-rental-protos/gen/service/trip"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/trip-service/internal/adapter/grpc/dto"
	"carsharing/trip-service/internal/model"
)

type TripService interface {
	StartTrip(ctx context.Context, bookingID string) (string, error)
	GetTrip(ctx context.Context, id string) (model.Trip, error)
	ListTrips(ctx context.Context, filter model.TripFilter) ([]model.Trip, error)
	EndTrip(ctx context.Context, id string) error
	CancelTrip(ctx context.Context, id string, reason *string) error
	GetTripSummary(ctx context.Context, tripID string) (model.TripSummary, error)
	GetTripStatusHistory(ctx context.Context, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error)
}

type TripHandler struct {
	tripsvc.UnimplementedTripServiceServer
	log     *slog.Logger
	service TripService
}

func NewTripHandler(log *slog.Logger, service TripService) *TripHandler {
	return &TripHandler{
		log:     pkglog.WithComponent(log, "handler.TripHandler"),
		service: service,
	}
}

func (h *TripHandler) StartTrip(ctx context.Context, req *tripsvc.StartTripRequest) (*tripsvc.StartTripResponse, error) {
	id, err := h.service.StartTrip(ctx, req.BookingId)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}
	return &tripsvc.StartTripResponse{Id: id}, nil
}

func (h *TripHandler) GetTrip(ctx context.Context, req *tripsvc.GetTripRequest) (*tripsvc.GetTripResponse, error) {
	trip, err := h.service.GetTrip(ctx, req.Id)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}
	return &tripsvc.GetTripResponse{Trip: dto.TripToProto(trip)}, nil
}

func (h *TripHandler) ListTrips(ctx context.Context, req *tripsvc.ListTripsRequest) (*tripsvc.ListTripsResponse, error) {
	trips, err := h.service.ListTrips(ctx, dto.FilterFromProto(req))
	if err != nil {
		return nil, dto.ToStatusError(err)
	}
	protos := make([]*basetripmb.Trip, len(trips))
	for i, t := range trips {
		protos[i] = dto.TripToProto(t)
	}
	return &tripsvc.ListTripsResponse{Trips: protos}, nil
}

func (h *TripHandler) EndTrip(ctx context.Context, req *tripsvc.EndTripRequest) (*emptypb.Empty, error) {
	if err := h.service.EndTrip(ctx, req.Id); err != nil {
		return nil, dto.ToStatusError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *TripHandler) CancelTrip(ctx context.Context, req *tripsvc.CancelTripRequest) (*emptypb.Empty, error) {
	if err := h.service.CancelTrip(ctx, req.Id, req.Reason); err != nil {
		return nil, dto.ToStatusError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *TripHandler) GetTripSummary(ctx context.Context, req *tripsvc.GetTripSummaryRequest) (*tripsvc.GetTripSummaryResponse, error) {
	summary, err := h.service.GetTripSummary(ctx, req.Id)
	if err != nil {
		return nil, dto.ToStatusError(err)
	}
	return &tripsvc.GetTripSummaryResponse{Summary: dto.TripSummaryToProto(summary)}, nil
}

func (h *TripHandler) GetTripStatusHistory(ctx context.Context, req *tripsvc.GetTripStatusHistoryRequest) (*tripsvc.GetTripStatusHistoryResponse, error) {
	history, err := h.service.GetTripStatusHistory(ctx, dto.StatusHistoryFilterFromProto(req))
	if err != nil {
		return nil, dto.ToStatusError(err)
	}
	protos := make([]*basetripmb.TripStatusReading, len(history))
	for i, r := range history {
		protos[i] = dto.TripStatusReadingToProto(r)
	}
	return &tripsvc.GetTripStatusHistoryResponse{StatusHistory: protos}, nil
}
