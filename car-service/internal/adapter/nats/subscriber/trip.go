package subscriber

import (
	"context"
	"log/slog"

	natsdto "carsharing/car-service/internal/adapter/nats/dto"
	pkglog "carsharing/shared/pkg/log"

	"github.com/nats-io/nats.go"
	eventtrip "github.com/sorawaslocked/car-rental-protos/gen/event/trip"
	"google.golang.org/protobuf/proto"
)

const (
	subjectTripStarted   = "trip.started"
	subjectTripEnded     = "trip.ended"
	subjectTripCancelled = "trip.cancelled"
)

type TripSubscriber struct {
	log          *slog.Logger
	conn         *nats.Conn
	eventHandler TripEventHandler
}

func NewTripSubscriber(log *slog.Logger, conn *nats.Conn, eventHandler TripEventHandler) *TripSubscriber {
	return &TripSubscriber{
		log:          pkglog.WithComponent(log, "adapter.nats.subscriber.TripSubscriber"),
		conn:         conn,
		eventHandler: eventHandler,
	}
}

func (s *TripSubscriber) Subscribe() error {
	subs := []struct {
		subject string
		handler nats.MsgHandler
	}{
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

func (s *TripSubscriber) handleTripStarted(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleTripStarted")

	var pb eventtrip.TripStartedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		log.Error("unmarshalling trip started event", pkglog.Err(err))
		return
	}

	event := natsdto.TripStartedEventFromProto(&pb)
	if err := s.eventHandler.OnTripStarted(ctx, event); err != nil {
		log.Error("handling trip started event",
			slog.String("tripID", event.TripID),
			slog.String("carID", event.CarID),
			pkglog.Err(err),
		)
	}
}

func (s *TripSubscriber) handleTripEnded(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleTripEnded")

	var pb eventtrip.TripEndedEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		log.Error("unmarshalling trip ended event", pkglog.Err(err))
		return
	}

	event := natsdto.TripEndedEventFromProto(&pb)
	if err := s.eventHandler.OnTripEnded(ctx, event); err != nil {
		log.Error("handling trip ended event",
			slog.String("tripID", event.TripID),
			slog.String("carID", event.CarID),
			pkglog.Err(err),
		)
	}
}

func (s *TripSubscriber) handleTripCancelled(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleTripCancelled")

	var pb eventtrip.TripCancelledEvent
	if err := proto.Unmarshal(msg.Data, &pb); err != nil {
		log.Error("unmarshalling trip cancelled event", pkglog.Err(err))
		return
	}

	event := natsdto.TripCancelledEventFromProto(&pb)
	if err := s.eventHandler.OnTripCancelled(ctx, event); err != nil {
		log.Error("handling trip cancelled event",
			slog.String("tripID", event.TripID),
			slog.String("carID", event.CarID),
			pkglog.Err(err),
		)
	}
}
