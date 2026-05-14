package client

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"

	basebookingpb "github.com/sorawaslocked/car-rental-protos/gen/base/booking"
	bookingsvc "github.com/sorawaslocked/car-rental-protos/gen/service/booking"

	"github.com/sorawaslocked/car-rental-trip-service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/car-rental-trip-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-trip-service/internal/pkg/log"
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
	resp, err := c.client.GetBooking(ctx, &bookingsvc.GetBookingRequest{Id: id})
	if err != nil {
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
