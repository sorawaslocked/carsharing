package client

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	basebookingpb "carsharing/protos/gen/base/booking"
	bookingsvc "carsharing/protos/gen/service/booking"

	pkglog "carsharing/shared/pkg/log"
	pkgutils "carsharing/shared/pkg/utils"
	"carsharing/trip-service/internal/adapter/grpc/dto"
	"carsharing/trip-service/internal/model"
)

type BookingClient struct {
	log    *slog.Logger
	client bookingsvc.BookingServiceClient
}

func NewBookingClient(log *slog.Logger, conn *grpc.ClientConn) *BookingClient {
	return &BookingClient{
		log:    pkglog.WithComponent(log, "client.BookingClient"),
		client: bookingsvc.NewBookingServiceClient(conn),
	}
}

func (c *BookingClient) GetBooking(ctx context.Context, id string) (model.Booking, error) {
	log := pkglog.WithMethod(c.log, "GetBooking")
	log = pkglog.WithMetadata(log, pkgutils.MetadataFromCtx(ctx))

	resp, err := c.client.GetBooking(ctx, &bookingsvc.GetBookingRequest{Id: id})
	if err != nil {
		log.Error("failed to get booking", pkglog.Err(err))
		return model.Booking{}, err
	}
	return bookingFromProto(resp.Booking), nil
}

func bookingFromProto(b *basebookingpb.Booking) model.Booking {
	booking := model.Booking{
		ID:               b.Id,
		UserID:           b.UserId,
		CarID:            b.CarId,
		Status:           b.Status,
		CommittedPeriods: b.CommittedPeriods,
		PricingSnapshot:  dto.PricingSnapshotFromProto(b.PricingSnapshot),
	}
	if b.ExpiresAt != nil {
		booking.ExpiresAt = b.ExpiresAt.AsTime()
	}
	return booking
}
