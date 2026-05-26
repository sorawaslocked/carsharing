package subscriber

import (
	"context"
	"log/slog"

	natsdto "carsharing/api-gateway/internal/adapter/nats/dto"
	eventcarpb "carsharing/protos/gen/event/car"
	pkglog "carsharing/shared/pkg/log"
	nc "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const subjectCarStatusUpdated = "car.status.updated"

type CarSubscriber struct {
	conn    *nc.Conn
	handler CarStatusEventHandler
	log     *slog.Logger
	subs    []*nc.Subscription
}

func NewCarSubscriber(conn *nc.Conn, handler CarStatusEventHandler, log *slog.Logger) *CarSubscriber {
	return &CarSubscriber{
		conn:    conn,
		handler: handler,
		log:     pkglog.WithComponent(log, "nats.CarSubscriber"),
	}
}

func (s *CarSubscriber) Subscribe() error {
	log := pkglog.WithMethod(s.log, "Subscribe")

	sub, err := s.conn.Subscribe(subjectCarStatusUpdated, s.handleCarStatusUpdated)
	if err != nil {
		log.Error("subscribing to subject",
			slog.String("subject", subjectCarStatusUpdated),
			pkglog.Err(err),
		)

		return natsdto.ErrSubscribeFailed
	}

	s.subs = append(s.subs, sub)
	log.Info("subscribed", slog.String("subject", subjectCarStatusUpdated))

	return nil
}

func (s *CarSubscriber) Close() {
	log := pkglog.WithMethod(s.log, "Close")

	log.Info("draining nats connection")

	if err := s.conn.Drain(); err != nil {
		log.Error("draining nats connection", pkglog.Err(err))
		s.conn.Close()
	}

	s.subs = nil
}

func (s *CarSubscriber) handleCarStatusUpdated(msg *nc.Msg) {
	log := pkglog.WithMethod(s.log, "handleCarStatusUpdated")

	var event eventcarpb.CarStatusUpdatedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		log.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnCarStatusUpdated(context.Background(), natsdto.CarStatusUpdatedFromProto(&event)); err != nil {
		log.Error("handling event",
			slog.String("carID", event.GetCarId()),
			pkglog.Err(err),
		)
	}
}
