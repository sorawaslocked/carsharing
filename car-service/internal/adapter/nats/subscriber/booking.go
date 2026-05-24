package subscriber

import (
	"context"
	"log/slog"

	natsdto "carsharing/car-service/internal/adapter/nats/dto"
	pkglog "carsharing/shared/pkg/log"

	eventbooking "carsharing/protos/gen/event/booking"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	subjectBookingCreated   = "booking.created"
	subjectBookingCancelled = "booking.cancelled"
	subjectBookingExpired   = "booking.expired"
	subjectBookingCompleted = "booking.completed"
)

type BookingSubscriber struct {
	log          *slog.Logger
	conn         *nats.Conn
	eventHandler BookingEventHandler
}

func NewBookingSubscriber(log *slog.Logger, conn *nats.Conn, eventHandler BookingEventHandler) *BookingSubscriber {
	return &BookingSubscriber{
		log:          pkglog.WithComponent(log, "adapter.nats.subscriber.BookingSubscriber"),
		conn:         conn,
		eventHandler: eventHandler,
	}
}

func (s *BookingSubscriber) Subscribe() error {
	subs := []struct {
		subject string
		handler nats.MsgHandler
	}{
		{subjectBookingCreated, s.handleBookingCreated},
		{subjectBookingCancelled, s.handleBookingCancelled},
		{subjectBookingExpired, s.handleBookingExpired},
		{subjectBookingCompleted, s.handleBookingCompleted},
	}

	for _, sub := range subs {
		if _, err := s.conn.Subscribe(sub.subject, sub.handler); err != nil {
			return err
		}
		s.log.Info("subscribed to NATS subject", slog.String("subject", sub.subject))
	}

	return nil
}

func (s *BookingSubscriber) handleBookingCreated(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleBookingCreated")

	var pb eventbooking.BookingCreatedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		log.Error("unmarshalling booking created event", pkglog.Err(err))
		return
	}

	event := natsdto.BookingCreatedEventFromProto(&pb)
	if err := s.eventHandler.OnBookingCreated(ctx, event); err != nil {
		log.Error("handling booking created event",
			slog.String("bookingID", event.BookingID),
			slog.String("carID", event.CarID),
			pkglog.Err(err),
		)
	}
}

func (s *BookingSubscriber) handleBookingCancelled(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleBookingCancelled")

	var pb eventbooking.BookingCancelledEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		log.Error("unmarshalling booking cancelled event", pkglog.Err(err))
		return
	}

	event := natsdto.BookingCancelledEventFromProto(&pb)
	if err := s.eventHandler.OnBookingCancelled(ctx, event); err != nil {
		log.Error("handling booking cancelled event",
			slog.String("bookingID", event.BookingID),
			slog.String("carID", event.CarID),
			pkglog.Err(err),
		)
	}
}

func (s *BookingSubscriber) handleBookingExpired(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleBookingExpired")

	var pb eventbooking.BookingExpiredEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		log.Error("unmarshalling booking expired event", pkglog.Err(err))
		return
	}

	event := natsdto.BookingExpiredEventFromProto(&pb)
	if err := s.eventHandler.OnBookingExpired(ctx, event); err != nil {
		log.Error("handling booking expired event",
			slog.String("bookingID", event.BookingID),
			slog.String("carID", event.CarID),
			pkglog.Err(err),
		)
	}
}

func (s *BookingSubscriber) handleBookingCompleted(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleBookingCompleted")

	var pb eventbooking.BookingCompletedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		log.Error("unmarshalling booking completed event", pkglog.Err(err))
		return
	}

	event := natsdto.BookingCompletedEventFromProto(&pb)
	if err := s.eventHandler.OnBookingCompleted(ctx, event); err != nil {
		log.Error("handling booking completed event",
			slog.String("bookingID", event.BookingID),
			slog.String("carID", event.CarID),
			pkglog.Err(err),
		)
	}
}
