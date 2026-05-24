package publisher

import (
	"context"
	"log/slog"
	"time"

	"carsharing/booking-service/internal/model"
	eventbookingpb "carsharing/protos/gen/event/booking"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	subjectBookingCreated   = "booking.created"
	subjectBookingCancelled = "booking.cancelled"
	subjectBookingExpired   = "booking.expired"
	subjectBookingCompleted = "booking.completed"
)

type BookingPublisher struct {
	log  *slog.Logger
	conn *nats.Conn
}

func NewBookingPublisher(log *slog.Logger, conn *nats.Conn) *BookingPublisher {
	return &BookingPublisher{
		log:  pkglog.WithComponent(log, "adapter.nats.publisher.BookingPublisher"),
		conn: conn,
	}
}

func (p *BookingPublisher) PublishBookingCreated(ctx context.Context, booking model.Booking) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishBookingCreated"), utils.MetadataFromCtx(ctx))

	return p.publish(log, subjectBookingCreated, &eventbookingpb.BookingCreatedEvent{
		BookingId: booking.ID,
		CarId:     booking.CarID,
		UserId:    booking.UserID,
		StartsAt:  timestamppb.New(booking.CreatedAt),
		EndsAt:    timestamppb.New(booking.ExpiresAt),
	})
}

func (p *BookingPublisher) PublishBookingCancelled(ctx context.Context, booking model.Booking, reason string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishBookingCancelled"), utils.MetadataFromCtx(ctx))

	return p.publish(log, subjectBookingCancelled, &eventbookingpb.BookingCancelledEvent{
		BookingId: booking.ID,
		CarId:     booking.CarID,
		UserId:    booking.UserID,
		Reason:    reason,
	})
}

func (p *BookingPublisher) PublishBookingExpired(ctx context.Context, booking model.Booking) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishBookingExpired"), utils.MetadataFromCtx(ctx))

	return p.publish(log, subjectBookingExpired, &eventbookingpb.BookingExpiredEvent{
		BookingId: booking.ID,
		CarId:     booking.CarID,
		UserId:    booking.UserID,
		ExpiredAt: timestamppb.New(time.Now()),
	})
}

func (p *BookingPublisher) PublishBookingCompleted(ctx context.Context, booking model.Booking) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(p.log, "PublishBookingCompleted"), utils.MetadataFromCtx(ctx))

	return p.publish(log, subjectBookingCompleted, &eventbookingpb.BookingCompletedEvent{
		BookingId:   booking.ID,
		CarId:       booking.CarID,
		UserId:      booking.UserID,
		CompletedAt: timestamppb.New(time.Now()),
	})
}

func (p *BookingPublisher) publish(log *slog.Logger, subject string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("marshalling event", slog.String("subject", subject), pkglog.Err(err))
		return model.ErrNats
	}

	if err := p.conn.Publish(subject, data); err != nil {
		log.Error("publishing event", slog.String("subject", subject), pkglog.Err(err))
		return model.ErrNats
	}

	return nil
}
