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

func NewCarSubscriber(conn *nc.Conn, handler CarStatusEventHandler, logger *slog.Logger) *CarSubscriber {
	return &CarSubscriber{
		conn:    conn,
		handler: handler,
		log:     pkglog.WithComponent(logger, "nats.CarSubscriber"),
	}
}

func (s *CarSubscriber) Subscribe() error {
	logger := pkglog.WithMethod(s.log, "Subscribe")

	sub, err := s.conn.Subscribe(subjectCarStatusUpdated, s.handleCarStatusUpdated)
	if err != nil {
		logger.Error("subscribing to subject",
			slog.String("subject", subjectCarStatusUpdated),
			pkglog.Err(err),
		)

		return natsdto.ErrSubscribeFailed
	}

	s.subs = append(s.subs, sub)
	logger.Info("subscribed", slog.String("subject", subjectCarStatusUpdated))

	return nil
}

func (s *CarSubscriber) Close() {
	logger := pkglog.WithMethod(s.log, "Close")

	logger.Info("draining nats connection")

	if err := s.conn.Drain(); err != nil {
		logger.Error("draining nats connection", pkglog.Err(err))
		s.conn.Close()
	}

	s.subs = nil
}

func (s *CarSubscriber) handleCarStatusUpdated(msg *nc.Msg) {
	logger := pkglog.WithMethod(s.log, "handleCarStatusUpdated")

	var event eventcarpb.CarStatusUpdatedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		logger.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnCarStatusUpdated(context.Background(), natsdto.CarStatusUpdatedFromProto(&event)); err != nil {
		logger.Error("handling event",
			slog.String("carID", event.GetCarId()),
			pkglog.Err(err),
		)
	}
}
