package nats

import (
	"context"
	"log/slog"

	natsio "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventtripmb "github.com/sorawaslocked/car-rental-protos/gen/event/trip"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/trip-service/internal/model"
)

const (
	subjectTripStarted   = "trip.started"
	subjectTripEnded     = "trip.ended"
	subjectTripCancelled = "trip.cancelled"
)

type Publisher struct {
	log  *slog.Logger
	conn *natsio.Conn
}

func NewPublisher(log *slog.Logger, conn *natsio.Conn) *Publisher {
	return &Publisher{
		log:  pkglog.WithComponent(log, "nats.Publisher"),
		conn: conn,
	}
}

func (p *Publisher) PublishTripStarted(_ context.Context, trip model.Trip) error {
	return p.publish(subjectTripStarted, &eventtripmb.TripStartedEvent{
		TripId:    trip.ID,
		BookingId: trip.BookingID,
		CarId:     trip.CarID,
		UserId:    trip.UserID,
		StartedAt: timestamppb.New(trip.StartedAt),
	})
}

func (p *Publisher) PublishTripEnded(_ context.Context, trip model.Trip) error {
	var endedAt *timestamppb.Timestamp
	if trip.EndedAt != nil {
		endedAt = timestamppb.New(*trip.EndedAt)
	}
	return p.publish(subjectTripEnded, &eventtripmb.TripEndedEvent{
		TripId:    trip.ID,
		BookingId: trip.BookingID,
		CarId:     trip.CarID,
		UserId:    trip.UserID,
		EndedAt:   endedAt,
	})
}

func (p *Publisher) PublishTripCancelled(_ context.Context, trip model.Trip) error {
	reason := ""
	if trip.CancelReason != nil {
		reason = *trip.CancelReason
	}
	return p.publish(subjectTripCancelled, &eventtripmb.TripCancelledEvent{
		TripId:      trip.ID,
		BookingId:   trip.BookingID,
		CarId:       trip.CarID,
		UserId:      trip.UserID,
		Reason:      reason,
		CancelledAt: timestamppb.New(trip.UpdatedAt),
	})
}

func (p *Publisher) publish(subject string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		p.log.Error("failed to marshal event", pkglog.Err(err), slog.String("subject", subject))
		return model.ErrNATS
	}
	if err = p.conn.Publish(subject, data); err != nil {
		p.log.Error("failed to publish event", pkglog.Err(err), slog.String("subject", subject))
		return model.ErrNATS
	}
	return nil
}
