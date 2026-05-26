package subscriber

import (
	"context"
	"log/slog"

	natsdto "carsharing/api-gateway/internal/adapter/nats/dto"
	eventuserpb "carsharing/protos/gen/event/user"
	pkglog "carsharing/shared/pkg/log"
	nc "github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	subjectUserCreated = "user.created"
	subjectUserUpdated = "user.updated"
	subjectUserDeleted = "user.deleted"
)

type UserSubscriber struct {
	conn    *nc.Conn
	handler UserEventHandler
	log     *slog.Logger
	subs    []*nc.Subscription
}

func NewUserSubscriber(conn *nc.Conn, handler UserEventHandler, log *slog.Logger) *UserSubscriber {
	return &UserSubscriber{
		conn:    conn,
		handler: handler,
		log:     pkglog.WithComponent(log, "nats.UserSubscriber"),
	}
}

func (s *UserSubscriber) Subscribe() error {
	type entry struct {
		subject string
		handler nc.MsgHandler
	}

	entries := []entry{
		{subjectUserCreated, s.handleUserCreated},
		{subjectUserUpdated, s.handleUserUpdated},
		{subjectUserDeleted, s.handleUserDeleted},
	}

	log := pkglog.WithMethod(s.log, "Subscribe")

	for _, e := range entries {
		sub, err := s.conn.Subscribe(e.subject, e.handler)
		if err != nil {
			log.Error("subscribing to subject",
				slog.String("subject", e.subject),
				pkglog.Err(err),
			)

			return natsdto.ErrSubscribeFailed
		}

		s.subs = append(s.subs, sub)
		log.Info("subscribed", slog.String("subject", e.subject))
	}

	return nil
}

func (s *UserSubscriber) Close() {
	log := pkglog.WithMethod(s.log, "Close")

	log.Info("draining nats connection")

	if err := s.conn.Drain(); err != nil {
		log.Error("draining nats connection", pkglog.Err(err))
		s.conn.Close()
	}

	s.subs = nil
}

func (s *UserSubscriber) handleUserCreated(msg *nc.Msg) {
	log := pkglog.WithMethod(s.log, "handleUserCreated")

	var event eventuserpb.UserCreatedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		log.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnUserCreated(context.Background(), event.GetId()); err != nil {
		log.Error("handling event",
			slog.String("userID", event.GetId()),
			pkglog.Err(err),
		)
	}
}

func (s *UserSubscriber) handleUserUpdated(msg *nc.Msg) {
	log := pkglog.WithMethod(s.log, "handleUserUpdated")

	var event eventuserpb.UserUpdatedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		log.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnUserUpdated(context.Background(), event.GetId(), event.GetIsSecurityUpdate()); err != nil {
		log.Error("handling event",
			slog.String("userID", event.GetId()),
			pkglog.Err(err),
		)
	}
}

func (s *UserSubscriber) handleUserDeleted(msg *nc.Msg) {
	log := pkglog.WithMethod(s.log, "handleUserDeleted")

	var event eventuserpb.UserDeletedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		log.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnUserDeleted(context.Background(), event.GetId()); err != nil {
		log.Error("handling event",
			slog.String("userID", event.GetId()),
			pkglog.Err(err),
		)
	}
}
