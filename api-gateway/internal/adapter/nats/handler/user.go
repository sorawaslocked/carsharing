package handler

import (
	"context"
	"log/slog"

	nc "github.com/nats-io/nats.go"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/nats/dto"
	pkglog "github.com/sorawaslocked/car-rental-api-gateway/internal/pkg/log"
	eventuserpb "github.com/sorawaslocked/car-rental-protos/gen/event/user"
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

func NewUserSubscriber(conn *nc.Conn, handler UserEventHandler, logger *slog.Logger) *UserSubscriber {
	return &UserSubscriber{
		conn:    conn,
		handler: handler,
		log:     pkglog.WithComponent(logger, "nats.UserSubscriber"),
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

	for _, e := range entries {
		sub, err := s.conn.Subscribe(e.subject, e.handler)
		if err != nil {
			logger := pkglog.WithMethod(s.log, "Subscribe")
			logger.Error("subscribing to subject",
				slog.String("subject", e.subject),
				pkglog.Err(err),
			)

			return dto.ErrSubscribeFailed
		}

		s.subs = append(s.subs, sub)
		s.log.Info("subscribed", slog.String("subject", e.subject))
	}

	return nil
}

func (s *UserSubscriber) Close() {
	logger := pkglog.WithMethod(s.log, "Close")

	logger.Info("draining nats connection")

	if err := s.conn.Drain(); err != nil {
		logger.Error("draining nats connection", pkglog.Err(err))
		s.conn.Close()
	}

	s.subs = nil
}

func (s *UserSubscriber) handleUserCreated(msg *nc.Msg) {
	logger := pkglog.WithMethod(s.log, "handleUserCreated")

	var event eventuserpb.CreateEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		logger.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnUserCreated(context.Background(), event.GetID()); err != nil {
		logger.Error("handling event",
			slog.String("userID", event.GetID()),
			pkglog.Err(err),
		)
	}
}

func (s *UserSubscriber) handleUserUpdated(msg *nc.Msg) {
	logger := pkglog.WithMethod(s.log, "handleUserUpdated")

	var event eventuserpb.UpdateEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		logger.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnUserUpdated(context.Background(), event.GetID(), event.GetIsSecurityUpdate()); err != nil {
		logger.Error("handling event",
			slog.String("userID", event.GetID()),
			pkglog.Err(err),
		)
	}
}

func (s *UserSubscriber) handleUserDeleted(msg *nc.Msg) {
	logger := pkglog.WithMethod(s.log, "handleUserDeleted")

	var event eventuserpb.DeleteEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		logger.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnUserDeleted(context.Background(), event.GetID()); err != nil {
		logger.Error("handling event",
			slog.String("userID", event.GetID()),
			pkglog.Err(err),
		)
	}
}
