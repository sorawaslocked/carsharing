package handler

import (
	"context"
	"log/slog"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/utils"
	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	bookingsvc "github.com/sorawaslocked/car-rental-protos/gen/service/booking"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookingHandler struct {
	client bookingsvc.BookingServiceClient
	log    *slog.Logger
}

func NewBookingHandler(client bookingsvc.BookingServiceClient, logger *slog.Logger) *BookingHandler {
	return &BookingHandler{
		client: client,
		log:    pkglog.WithComponent(logger, "grpc.BookingHandler"),
	}
}

func (h *BookingHandler) Create(ctx context.Context, data model.BookingCreate) (string, error) {
	logger := pkglog.WithMethod(h.log, "Create")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &bookingsvc.CreateBookingRequest{
		UserId:        data.UserID,
		CarId:         data.CarID,
		PricingRuleId: data.PricingRuleID,
	}
	if data.CommittedPeriods != nil {
		req.CommittedPeriods = data.CommittedPeriods
	}

	res, err := h.client.CreateBooking(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *BookingHandler) Get(ctx context.Context, id string) (model.Booking, error) {
	logger := pkglog.WithMethod(h.log, "Get")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	res, err := h.client.GetBooking(ctx, &bookingsvc.GetBookingRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return model.Booking{}, dto.FromGrpcErr(err)
	}

	return dto.BookingFromProto(res.GetBooking()), nil
}

func (h *BookingHandler) List(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error) {
	logger := pkglog.WithMethod(h.log, "List")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &bookingsvc.ListBookingsRequest{
		UserId:        filter.UserID,
		CarId:         filter.CarID,
		Status:        filter.Status,
		PricingRuleId: filter.PricingRuleID,
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.ListBookings(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	bookings := make([]model.Booking, len(res.GetBookings()))
	for i, b := range res.GetBookings() {
		bookings[i] = dto.BookingFromProto(b)
	}

	return bookings, nil
}

func (h *BookingHandler) Cancel(ctx context.Context, id string) error {
	logger := pkglog.WithMethod(h.log, "Cancel")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := h.client.CancelBooking(ctx, &bookingsvc.CancelBookingRequest{Id: id})
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *BookingHandler) UpdateStatus(ctx context.Context, id string, data model.BookingStatusUpdate) error {
	logger := pkglog.WithMethod(h.log, "UpdateStatus")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &bookingsvc.UpdateBookingStatusRequest{
		Id:     id,
		Status: data.Status,
		Reason: data.Reason,
	}

	_, err := h.client.UpdateBookingStatus(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *BookingHandler) GetStatusHistory(ctx context.Context, id string, filter model.BookingStatusReadingFilter) ([]model.BookingStatusReading, error) {
	logger := pkglog.WithMethod(h.log, "GetStatusHistory")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	req := &bookingsvc.GetBookingStatusHistoryRequest{Id: id}
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

	res, err := h.client.GetBookingStatusHistory(ctx, req)
	if err != nil {
		if dto.IsSystemErr(err) {
			logger.Error("grpc call failed", pkglog.Err(err))
		}

		return nil, dto.FromGrpcErr(err)
	}

	readings := make([]model.BookingStatusReading, len(res.GetStatusHistory()))
	for i, r := range res.GetStatusHistory() {
		readings[i] = dto.BookingStatusReadingFromProto(r)
	}

	return readings, nil
}
