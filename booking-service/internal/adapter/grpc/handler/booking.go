package handler

import (
	"context"
	"log/slog"

	"carsharing/booking-service/internal/adapter/grpc/dto"
	"carsharing/booking-service/internal/model"
	basebookingpb "carsharing/protos/gen/base/booking"
	servicebookingpb "carsharing/protos/gen/service/booking"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BookingHandler struct {
	servicebookingpb.UnimplementedBookingServiceServer
	log *slog.Logger
	svc BookingService
}

func NewBookingHandler(log *slog.Logger, svc BookingService) *BookingHandler {
	return &BookingHandler{
		log: pkglog.WithComponent(log, "grpc.BookingHandler"),
		svc: svc,
	}
}

func (h *BookingHandler) CreateBooking(ctx context.Context, req *servicebookingpb.CreateBookingRequest) (*servicebookingpb.CreateBookingResponse, error) {
	log := pkglog.WithMethod(h.log, "CreateBooking")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	data := model.BookingCreate{
		UserID:           req.UserId,
		CarID:            req.CarId,
		PricingRuleID:    req.PricingRuleId,
		CommittedPeriods: req.CommittedPeriods,
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		return nil, dto.ToGRPCError(err)
	}

	log.Info("booking created", slog.String("id", id))

	return &servicebookingpb.CreateBookingResponse{Id: id}, nil
}

func (h *BookingHandler) GetBooking(ctx context.Context, req *servicebookingpb.GetBookingRequest) (*servicebookingpb.GetBookingResponse, error) {
	log := pkglog.WithMethod(h.log, "GetBooking")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	booking, err := h.svc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, dto.ToGRPCError(err)
	}

	log.Info("booking retrieved", slog.String("id", req.Id))

	return &servicebookingpb.GetBookingResponse{Booking: dto.BookingToProto(booking)}, nil
}

func (h *BookingHandler) ListBookings(ctx context.Context, req *servicebookingpb.ListBookingsRequest) (*servicebookingpb.ListBookingsResponse, error) {
	log := pkglog.WithMethod(h.log, "ListBookings")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	bookings, err := h.svc.List(ctx, dto.BookingListFilterFromProto(req))
	if err != nil {
		return nil, dto.ToGRPCError(err)
	}

	log.Info("bookings listed", slog.Int("count", len(bookings)))

	pbBookings := make([]*basebookingpb.Booking, 0, len(bookings))
	for _, b := range bookings {
		pbBookings = append(pbBookings, dto.BookingToProto(b))
	}

	return &servicebookingpb.ListBookingsResponse{Bookings: pbBookings}, nil
}

func (h *BookingHandler) CancelBooking(ctx context.Context, req *servicebookingpb.CancelBookingRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMethod(h.log, "CancelBooking")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	if err := h.svc.Cancel(ctx, req.Id, nil); err != nil {
		return nil, dto.ToGRPCError(err)
	}

	log.Info("booking cancelled", slog.String("id", req.Id))

	return &emptypb.Empty{}, nil
}

func (h *BookingHandler) UpdateBookingStatus(ctx context.Context, req *servicebookingpb.UpdateBookingStatusRequest) (*emptypb.Empty, error) {
	log := pkglog.WithMethod(h.log, "UpdateBookingStatus")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	if err := h.svc.UpdateStatus(ctx, req.Id, req.Status, req.Reason); err != nil {
		return nil, dto.ToGRPCError(err)
	}

	log.Info("booking status updated", slog.String("id", req.Id), slog.String("status", req.Status))

	return &emptypb.Empty{}, nil
}

func (h *BookingHandler) GetBookingStatusHistory(ctx context.Context, req *servicebookingpb.GetBookingStatusHistoryRequest) (*servicebookingpb.GetBookingStatusHistoryResponse, error) {
	log := pkglog.WithMethod(h.log, "GetBookingStatusHistory")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	history, err := h.svc.GetStatusHistory(ctx, dto.BookingStatusHistoryFilterFromProto(req))
	if err != nil {
		return nil, dto.ToGRPCError(err)
	}

	log.Info("status history retrieved", slog.String("bookingID", req.Id), slog.Int("count", len(history)))

	pbHistory := make([]*basebookingpb.BookingStatusReading, 0, len(history))
	for _, r := range history {
		pbHistory = append(pbHistory, dto.BookingStatusReadingToProto(r))
	}

	return &servicebookingpb.GetBookingStatusHistoryResponse{StatusHistory: pbHistory}, nil
}
