package nats

import (
	"context"
	"log/slog"

	"carsharing/car-service/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"github.com/nats-io/nats.go"
	eventbooking "github.com/sorawaslocked/car-rental-protos/gen/event/booking"
	eventtrip "github.com/sorawaslocked/car-rental-protos/gen/event/trip"
	"google.golang.org/protobuf/proto"
)

const (
	subjectBookingCreated   = "booking.created"
	subjectBookingCancelled = "booking.cancelled"
	subjectBookingExpired   = "booking.expired"
	subjectBookingCompleted = "booking.completed"
	subjectTripStarted      = "trip.started"
	subjectTripEnded        = "trip.ended"
	subjectTripCancelled    = "trip.cancelled"
)

type CarEventHandler interface {
	OnBookingCreated(ctx context.Context, event model.BookingCreatedEvent) error
	OnBookingCancelled(ctx context.Context, event model.BookingCancelledEvent) error
	OnBookingExpired(ctx context.Context, event model.BookingExpiredEvent) error
	OnBookingCompleted(ctx context.Context, event model.BookingCompletedEvent) error
	OnTripStarted(ctx context.Context, event model.TripStartedEvent) error
	OnTripEnded(ctx context.Context, event model.TripEndedEvent) error
	OnTripCancelled(ctx context.Context, event model.TripCancelledEvent) error
}

type Subscriber struct {
	conn       *nats.Conn
	carHandler CarEventHandler
	log        *slog.Logger
}

func NewSubscriber(conn *nats.Conn, carHandler CarEventHandler, log *slog.Logger) *Subscriber {
	return &Subscriber{
		conn:       conn,
		carHandler: carHandler,
		log:        pkglog.WithComponent(log, "adapter.NATSSubscriber"),
	}
}

func (s *Subscriber) Subscribe() error {
	subs := []struct {
		subject string
		handler nats.MsgHandler
	}{
		{subjectBookingCreated, s.handleBookingCreated},
		{subjectBookingCancelled, s.handleBookingCancelled},
		{subjectBookingExpired, s.handleBookingExpired},
		{subjectBookingCompleted, s.handleBookingCompleted},
		{subjectTripStarted, s.handleTripStarted},
		{subjectTripEnded, s.handleTripEnded},
		{subjectTripCancelled, s.handleTripCancelled},
	}

	for _, sub := range subs {
		if _, err := s.conn.Subscribe(sub.subject, sub.handler); err != nil {
			return err
		}
		s.log.Info("subscribed to NATS subject", slog.String("subject", sub.subject))
	}

	return nil
}

func (s *Subscriber) handleBookingCreated(msg *nats.Msg) {
	var pb eventbooking.BookingCreatedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		s.log.Error("failed to unmarshal BookingCreatedEvent", pkglog.Err(err))
		return
	}

	event := model.BookingCreatedEvent{
		BookingID: pb.GetBookingId(),
		CarID:     pb.GetCarId(),
		UserID:    pb.GetUserId(),
	}
	if pb.GetStartsAt() != nil {
		event.StartsAt = pb.GetStartsAt().AsTime()
	}
	if pb.GetEndsAt() != nil {
		event.EndsAt = pb.GetEndsAt().AsTime()
	}

	if err := s.carHandler.OnBookingCreated(context.Background(), event); err != nil {
		s.log.Error("OnBookingCreated failed",
			pkglog.Err(err),
			slog.String("bookingID", event.BookingID),
			slog.String("carID", event.CarID),
		)
	}
}

func (s *Subscriber) handleBookingCancelled(msg *nats.Msg) {
	var pb eventbooking.BookingCancelledEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		s.log.Error("failed to unmarshal BookingCancelledEvent", pkglog.Err(err))
		return
	}

	event := model.BookingCancelledEvent{
		BookingID: pb.GetBookingId(),
		CarID:     pb.GetCarId(),
		UserID:    pb.GetUserId(),
		Reason:    pb.GetReason(),
	}

	if err := s.carHandler.OnBookingCancelled(context.Background(), event); err != nil {
		s.log.Error("OnBookingCancelled failed",
			pkglog.Err(err),
			slog.String("bookingID", event.BookingID),
			slog.String("carID", event.CarID),
		)
	}
}

func (s *Subscriber) handleTripStarted(msg *nats.Msg) {
	var pb eventtrip.TripStartedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		s.log.Error("failed to unmarshal TripStartedEvent", pkglog.Err(err))
		return
	}

	event := model.TripStartedEvent{
		TripID:    pb.GetTripId(),
		BookingID: pb.GetBookingId(),
		CarID:     pb.GetCarId(),
		UserID:    pb.GetUserId(),
	}
	if pb.GetStartedAt() != nil {
		event.StartedAt = pb.GetStartedAt().AsTime()
	}

	if err := s.carHandler.OnTripStarted(context.Background(), event); err != nil {
		s.log.Error("OnTripStarted failed",
			pkglog.Err(err),
			slog.String("tripID", event.TripID),
			slog.String("carID", event.CarID),
		)
	}
}

func (s *Subscriber) handleTripEnded(msg *nats.Msg) {
	var pb eventtrip.TripEndedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		s.log.Error("failed to unmarshal TripEndedEvent", pkglog.Err(err))
		return
	}

	event := model.TripEndedEvent{
		TripID:    pb.GetTripId(),
		BookingID: pb.GetBookingId(),
		CarID:     pb.GetCarId(),
		UserID:    pb.GetUserId(),
	}
	if pb.GetEndedAt() != nil {
		event.EndedAt = pb.GetEndedAt().AsTime()
	}

	if err := s.carHandler.OnTripEnded(context.Background(), event); err != nil {
		s.log.Error("OnTripEnded failed",
			pkglog.Err(err),
			slog.String("tripID", event.TripID),
			slog.String("carID", event.CarID),
		)
	}
}

func (s *Subscriber) handleBookingExpired(msg *nats.Msg) {
	var pb eventbooking.BookingExpiredEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		s.log.Error("failed to unmarshal BookingExpiredEvent", pkglog.Err(err))
		return
	}

	event := model.BookingExpiredEvent{
		BookingID: pb.GetBookingId(),
		CarID:     pb.GetCarId(),
		UserID:    pb.GetUserId(),
	}
	if pb.GetExpiredAt() != nil {
		event.ExpiredAt = pb.GetExpiredAt().AsTime()
	}

	if err := s.carHandler.OnBookingExpired(context.Background(), event); err != nil {
		s.log.Error("OnBookingExpired failed",
			pkglog.Err(err),
			slog.String("bookingID", event.BookingID),
			slog.String("carID", event.CarID),
		)
	}
}

func (s *Subscriber) handleBookingCompleted(msg *nats.Msg) {
	var pb eventbooking.BookingCompletedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		s.log.Error("failed to unmarshal BookingCompletedEvent", pkglog.Err(err))
		return
	}

	event := model.BookingCompletedEvent{
		BookingID: pb.GetBookingId(),
		CarID:     pb.GetCarId(),
		UserID:    pb.GetUserId(),
	}
	if pb.GetCompletedAt() != nil {
		event.CompletedAt = pb.GetCompletedAt().AsTime()
	}

	if err := s.carHandler.OnBookingCompleted(context.Background(), event); err != nil {
		s.log.Error("OnBookingCompleted failed",
			pkglog.Err(err),
			slog.String("bookingID", event.BookingID),
			slog.String("carID", event.CarID),
		)
	}
}

func (s *Subscriber) handleTripCancelled(msg *nats.Msg) {
	var pb eventtrip.TripCancelledEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		s.log.Error("failed to unmarshal TripCancelledEvent", pkglog.Err(err))
		return
	}

	event := model.TripCancelledEvent{
		TripID:    pb.GetTripId(),
		BookingID: pb.GetBookingId(),
		CarID:     pb.GetCarId(),
		UserID:    pb.GetUserId(),
		Reason:    pb.GetReason(),
	}
	if pb.GetCancelledAt() != nil {
		event.CancelledAt = pb.GetCancelledAt().AsTime()
	}

	if err := s.carHandler.OnTripCancelled(context.Background(), event); err != nil {
		s.log.Error("OnTripCancelled failed",
			pkglog.Err(err),
			slog.String("tripID", event.TripID),
			slog.String("carID", event.CarID),
		)
	}
}
