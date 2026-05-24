package handler

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	tripsvc "carsharing/protos/gen/service/trip"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TripHandler struct {
	client       tripsvc.TripServiceClient
	streamClient tripsvc.TripStreamServiceClient
	log          *slog.Logger
}

func NewTripHandler(client tripsvc.TripServiceClient, streamClient tripsvc.TripStreamServiceClient, logger *slog.Logger) *TripHandler {
	return &TripHandler{
		client:       client,
		streamClient: streamClient,
		log:          pkglog.WithComponent(logger, "grpc.TripHandler"),
	}
}

func (h *TripHandler) Start(ctx context.Context, bookingID string) (string, error) {
	logger := pkglog.WithMethod(h.log, "Start")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.StartTrip(ctx, &tripsvc.StartTripRequest{BookingId: bookingID})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *TripHandler) Get(ctx context.Context, id string) (model.Trip, error) {
	logger := pkglog.WithMethod(h.log, "Get")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetTrip(ctx, &tripsvc.GetTripRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.Trip{}, dto.FromGrpcErr(err)
	}

	return dto.TripFromProto(res.GetTrip()), nil
}

func (h *TripHandler) List(ctx context.Context, filter model.TripFilter) ([]model.Trip, error) {
	logger := pkglog.WithMethod(h.log, "List")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &tripsvc.ListTripsRequest{
		UserId: filter.UserID,
		CarId:  filter.CarID,
		Status: filter.Status,
	}
	if filter.StartedAfter != nil {
		req.StartedAfter = timestamppb.New(*filter.StartedAfter)
	}
	if filter.StartedBefore != nil {
		req.StartedBefore = timestamppb.New(*filter.StartedBefore)
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListTrips(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	trips := make([]model.Trip, len(res.GetTrips()))
	for i, t := range res.GetTrips() {
		trips[i] = dto.TripFromProto(t)
	}

	return trips, nil
}

func (h *TripHandler) End(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "End")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.EndTrip(ctx, &tripsvc.EndTripRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *TripHandler) Cancel(ctx context.Context, id string, reason *string) error {
	logger := pkglog.WithMethod(h.log, "Cancel")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.CancelTrip(ctx, &tripsvc.CancelTripRequest{Id: id, Reason: reason})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *TripHandler) GetSummary(ctx context.Context, id string) (model.TripSummary, error) {
	logger := pkglog.WithMethod(h.log, "GetSummary")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetTripSummary(ctx, &tripsvc.GetTripSummaryRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.TripSummary{}, dto.FromGrpcErr(err)
	}

	return dto.TripSummaryFromProto(res.GetSummary()), nil
}

func (h *TripHandler) GetStatusHistory(ctx context.Context, id string, filter model.TripStatusReadingFilter) ([]model.TripStatusReading, error) {
	logger := pkglog.WithMethod(h.log, "GetStatusHistory")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &tripsvc.GetTripStatusHistoryRequest{Id: id}
	if filter.From != nil {
		req.From = timestamppb.New(*filter.From)
	}
	if filter.To != nil {
		req.To = timestamppb.New(*filter.To)
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.GetTripStatusHistory(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	readings := make([]model.TripStatusReading, len(res.GetStatusHistory()))
	for i, r := range res.GetStatusHistory() {
		readings[i] = dto.TripStatusReadingFromProto(r)
	}

	return readings, nil
}
