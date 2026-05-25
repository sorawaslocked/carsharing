package handler

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/emptypb"

	basetripmb "carsharing/protos/gen/base/trip"
	tripsvc "carsharing/protos/gen/service/trip"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/adapter/grpc/dto"
)

type TripHandler struct {
	tripsvc.UnimplementedTripServiceServer
	log     *slog.Logger
	service TripService
}

func NewTripHandler(log *slog.Logger, service TripService) *TripHandler {
	return &TripHandler{
		log:     pkglog.WithComponent(log, "adapter.grpc.handler.TripHandler"),
		service: service,
	}
}

func (h *TripHandler) StartTrip(ctx context.Context, req *tripsvc.StartTripRequest) (*tripsvc.StartTripResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "StartTrip"), pkgutils.MetadataFromCtx(ctx))

	id, err := h.service.StartTrip(ctx, req.BookingId)
	if err != nil {
		log.Warn("starting trip", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}
	return &tripsvc.StartTripResponse{Id: id}, nil
}

func (h *TripHandler) GetTrip(ctx context.Context, req *tripsvc.GetTripRequest) (*tripsvc.GetTripResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetTrip"), pkgutils.MetadataFromCtx(ctx))

	trip, err := h.service.GetTrip(ctx, req.Id)
	if err != nil {
		log.Warn("getting trip", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}
	return &tripsvc.GetTripResponse{Trip: dto.TripToProto(trip)}, nil
}

func (h *TripHandler) ListTrips(ctx context.Context, req *tripsvc.ListTripsRequest) (*tripsvc.ListTripsResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "ListTrips"), pkgutils.MetadataFromCtx(ctx))

	trips, err := h.service.ListTrips(ctx, dto.FilterFromProto(req))
	if err != nil {
		log.Warn("listing trips", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}
	protos := make([]*basetripmb.Trip, len(trips))
	for i, t := range trips {
		protos[i] = dto.TripToProto(t)
	}
	return &tripsvc.ListTripsResponse{Trips: protos}, nil
}

func (h *TripHandler) EndTrip(ctx context.Context, req *tripsvc.EndTripRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "EndTrip"), pkgutils.MetadataFromCtx(ctx))

	if err := h.service.EndTrip(ctx, req.Id); err != nil {
		log.Warn("ending trip", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *TripHandler) CancelTrip(ctx context.Context, req *tripsvc.CancelTripRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CancelTrip"), pkgutils.MetadataFromCtx(ctx))

	if err := h.service.CancelTrip(ctx, req.Id, req.Reason); err != nil {
		log.Warn("cancelling trip", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}
	return &emptypb.Empty{}, nil
}

func (h *TripHandler) GetTripSummary(ctx context.Context, req *tripsvc.GetTripSummaryRequest) (*tripsvc.GetTripSummaryResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetTripSummary"), pkgutils.MetadataFromCtx(ctx))

	summary, err := h.service.GetTripSummary(ctx, req.Id)
	if err != nil {
		log.Warn("getting trip summary", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}
	return &tripsvc.GetTripSummaryResponse{Summary: dto.TripSummaryToProto(summary)}, nil
}

func (h *TripHandler) GetTripStatusHistory(ctx context.Context, req *tripsvc.GetTripStatusHistoryRequest) (*tripsvc.GetTripStatusHistoryResponse, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetTripStatusHistory"), pkgutils.MetadataFromCtx(ctx))

	history, err := h.service.GetTripStatusHistory(ctx, dto.StatusHistoryFilterFromProto(req))
	if err != nil {
		log.Warn("getting trip status history", pkglog.Err(err))
		return nil, dto.ToStatusError(err)
	}
	protos := make([]*basetripmb.TripStatusReading, len(history))
	for i, r := range history {
		protos[i] = dto.TripStatusReadingToProto(r)
	}
	return &tripsvc.GetTripStatusHistoryResponse{StatusHistory: protos}, nil
}
