package subscriber

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"

	eventtripb "carsharing/protos/gen/event/trip"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const subjectTripStarted = "trip.started"

type TripSubscriber struct {
	log     *slog.Logger
	conn    *nats.Conn
	handler TripEventHandler
}

func NewTripSubscriber(log *slog.Logger, conn *nats.Conn, handler TripEventHandler) *TripSubscriber {
	return &TripSubscriber{
		log:     pkglog.WithComponent(log, "adapter.nats.subscriber.TripSubscriber"),
		conn:    conn,
		handler: handler,
	}
}

func (s *TripSubscriber) Subscribe() error {
	_, err := s.conn.Subscribe(subjectTripStarted, s.handleTripStarted)
	return err
}

func (s *TripSubscriber) handleTripStarted(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleTripStarted")

	var event eventtripb.TripStartedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		log.Error("unmarshalling trip started event", pkglog.Err(err))
		return
	}

	if err := s.handler.Complete(ctx, event.BookingId); err != nil {
		log.Error("completing booking on trip.started",
			slog.String("bookingID", event.BookingId),
			slog.String("tripID", event.TripId),
			pkglog.Err(err),
		)
		return
	}

	log.Info("booking completed on trip.started",
		slog.String("bookingID", event.BookingId),
		slog.String("tripID", event.TripId),
	)
}
