package nats

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"

	natsgo "github.com/nats-io/nats.go"
	eventtripb "github.com/sorawaslocked/car-rental-protos/gen/event/trip"
	"google.golang.org/protobuf/proto"
)

const subjectTripStarted = "trip.started"

type TripEventHandler interface {
	Complete(ctx context.Context, bookingID string) error
}

type Consumer struct {
	log     *slog.Logger
	conn    *natsgo.Conn
	handler TripEventHandler
}

func NewConsumer(log *slog.Logger, conn *natsgo.Conn, handler TripEventHandler) *Consumer {
	return &Consumer{
		log:     pkglog.WithComponent(log, "nats.Consumer"),
		conn:    conn,
		handler: handler,
	}
}

func (c *Consumer) Subscribe(ctx context.Context) error {
	sub, err := c.conn.Subscribe(subjectTripStarted, func(msg *natsgo.Msg) {
		var event eventtripb.TripStartedEvent
		if err := proto.Unmarshal(msg.Data, &event); err != nil {
			c.log.Error("failed to unmarshal TripStartedEvent", pkglog.Err(err))
			return
		}

		if err := c.handler.Complete(context.Background(), event.BookingId); err != nil {
			c.log.Error("failed to complete booking on trip.started",
				slog.String("bookingID", event.BookingId),
				slog.String("tripID", event.TripId),
				pkglog.Err(err),
			)
			return
		}

		c.log.Info("booking completed on trip.started",
			slog.String("bookingID", event.BookingId),
			slog.String("tripID", event.TripId),
		)
	})
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
		c.log.Info("unsubscribed from trip.started")
	}()

	return nil
}
