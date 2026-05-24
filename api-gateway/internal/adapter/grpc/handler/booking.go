package handler

import (
	"context"
	"log/slog"

	"carsharing/api-gateway/internal/adapter/grpc/dto"
	"carsharing/api-gateway/internal/model"
	basepb "carsharing/protos/gen/base"
	bookingsvc "carsharing/protos/gen/service/booking"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
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
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(ctx))

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
		log.Warn("creating booking", pkglog.Err(err))

		return "", dto.FromGrpcErr(err)
	}

	return res.GetId(), nil
}

func (h *BookingHandler) Get(ctx context.Context, id string) (model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))

	res, err := h.client.GetBooking(ctx, &bookingsvc.GetBookingRequest{Id: id})
	if err != nil {
		log.Warn("getting booking", pkglog.Err(err))

		return model.Booking{}, dto.FromGrpcErr(err)
	}

	return dto.BookingFromProto(res.GetBooking()), nil
}

func (h *BookingHandler) List(ctx context.Context, filter model.BookingFilter) ([]model.Booking, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(ctx))

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
		log.Warn("listing bookings", pkglog.Err(err))

		return nil, dto.FromGrpcErr(err)
	}

	bookings := make([]model.Booking, len(res.GetBookings()))
	for i, b := range res.GetBookings() {
		bookings[i] = dto.BookingFromProto(b)
	}

	return bookings, nil
}

func (h *BookingHandler) Cancel(ctx context.Context, id string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Cancel"), utils.MetadataFromCtx(ctx))

	_, err := h.client.CancelBooking(ctx, &bookingsvc.CancelBookingRequest{Id: id})
	if err != nil {
		log.Warn("cancelling booking", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *BookingHandler) UpdateStatus(ctx context.Context, id string, data model.BookingStatusUpdate) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "UpdateStatus"), utils.MetadataFromCtx(ctx))

	_, err := h.client.UpdateBookingStatus(ctx, &bookingsvc.UpdateBookingStatusRequest{
		Id:     id,
		Status: data.Status,
		Reason: data.Reason,
	})
	if err != nil {
		log.Warn("updating booking status", pkglog.Err(err))

		return dto.FromGrpcErr(err)
	}

	return nil
}

func (h *BookingHandler) GetStatusHistory(ctx context.Context, id string, filter model.BookingStatusReadingFilter) ([]model.BookingStatusReading, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetStatusHistory"), utils.MetadataFromCtx(ctx))

	req := &bookingsvc.GetBookingStatusHistoryRequest{Id: id}
	if filter.From != nil || filter.To != nil {
		req.TimeRange = &basepb.TimeRange{}
		if filter.From != nil {
			req.TimeRange.From = timestamppb.New(*filter.From)
		}
		if filter.To != nil {
			req.TimeRange.To = timestamppb.New(*filter.To)
		}
	}
	if filter.Pagination != nil {
		req.Pagination = &basepb.Pagination{
			Limit:  filter.Pagination.Limit,
			Offset: filter.Pagination.Offset,
		}
	}

	res, err := h.client.GetBookingStatusHistory(ctx, req)
	if err != nil {
		log.Warn("getting booking status history", pkglog.Err(err))

		return nil, dto.FromGrpcErr(err)
	}

	readings := make([]model.BookingStatusReading, len(res.GetStatusHistory()))
	for i, r := range res.GetStatusHistory() {
		readings[i] = dto.BookingStatusReadingFromProto(r)
	}

	return readings, nil
}
