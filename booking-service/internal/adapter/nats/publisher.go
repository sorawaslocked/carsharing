package nats

import (
	"context"
	"log/slog"
	"time"

	"carsharing/booking-service/internal/model"
	eventbookingpb "carsharing/protos/gen/event/booking"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	natsgo "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	subjectBookingCreated   = "booking.created"
	subjectBookingCancelled = "booking.cancelled"
	subjectBookingExpired   = "booking.expired"
	subjectBookingCompleted = "booking.completed"
)

type Publisher struct {
	log  *slog.Logger
	conn *natsgo.Conn
}

func NewPublisher(log *slog.Logger, conn *natsgo.Conn) *Publisher {
	return &Publisher{
		log:  pkglog.WithComponent(log, "nats.Publisher"),
		conn: conn,
	}
}

func (p *Publisher) PublishBookingCreated(ctx context.Context, booking model.Booking) error {
	log := pkglog.WithMethod(p.log, "PublishBookingCreated")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	event := &eventbookingpb.BookingCreatedEvent{
		BookingId: booking.ID,
		CarId:     booking.CarID,
		UserId:    booking.UserID,
		StartsAt:  timestamppb.New(booking.CreatedAt),
		EndsAt:    timestamppb.New(booking.ExpiresAt),
	}

	return p.publish(log, subjectBookingCreated, event, booking.ID)
}

func (p *Publisher) PublishBookingCancelled(ctx context.Context, booking model.Booking, reason string) error {
	log := pkglog.WithMethod(p.log, "PublishBookingCancelled")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	event := &eventbookingpb.BookingCancelledEvent{
		BookingId: booking.ID,
		CarId:     booking.CarID,
		UserId:    booking.UserID,
		Reason:    reason,
	}

	return p.publish(log, subjectBookingCancelled, event, booking.ID)
}

func (p *Publisher) PublishBookingExpired(ctx context.Context, booking model.Booking) error {
	log := pkglog.WithMethod(p.log, "PublishBookingExpired")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	event := &eventbookingpb.BookingExpiredEvent{
		BookingId: booking.ID,
		CarId:     booking.CarID,
		UserId:    booking.UserID,
		ExpiredAt: timestamppb.New(time.Now()),
	}

	return p.publish(log, subjectBookingExpired, event, booking.ID)
}

func (p *Publisher) PublishBookingCompleted(ctx context.Context, booking model.Booking) error {
	log := pkglog.WithMethod(p.log, "PublishBookingCompleted")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	event := &eventbookingpb.BookingCompletedEvent{
		BookingId:   booking.ID,
		CarId:       booking.CarID,
		UserId:      booking.UserID,
		CompletedAt: timestamppb.New(time.Now()),
	}

	return p.publish(log, subjectBookingCompleted, event, booking.ID)
}

func (p *Publisher) publish(log *slog.Logger, subject string, msg proto.Message, bookingID string) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("failed to marshal event", slog.String("subject", subject), pkglog.Err(err))
		return err
	}

	if err := p.conn.Publish(subject, data); err != nil {
		log.Error("failed to publish event", slog.String("subject", subject), pkglog.Err(err))
		return err
	}

	log.Info("published event", slog.String("subject", subject), slog.String("bookingID", bookingID))
	return nil
}
