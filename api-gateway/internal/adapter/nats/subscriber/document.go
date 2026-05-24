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

const subjectDocumentAnalyzed = "user.document.analyzed"

type DocumentSubscriber struct {
	conn    *nc.Conn
	handler DocumentEventHandler
	log     *slog.Logger
	subs    []*nc.Subscription
}

func NewDocumentSubscriber(conn *nc.Conn, handler DocumentEventHandler, logger *slog.Logger) *DocumentSubscriber {
	return &DocumentSubscriber{
		conn:    conn,
		handler: handler,
		log:     pkglog.WithComponent(logger, "nats.DocumentSubscriber"),
	}
}

func (s *DocumentSubscriber) Subscribe() error {
	logger := pkglog.WithMethod(s.log, "Subscribe")

	sub, err := s.conn.Subscribe(subjectDocumentAnalyzed, s.handleDocumentAnalyzed)
	if err != nil {
		logger.Error("subscribing to subject",
			slog.String("subject", subjectDocumentAnalyzed),
			pkglog.Err(err),
		)

		return natsdto.ErrSubscribeFailed
	}

	s.subs = append(s.subs, sub)
	logger.Info("subscribed", slog.String("subject", subjectDocumentAnalyzed))

	return nil
}

func (s *DocumentSubscriber) Close() {
	logger := pkglog.WithMethod(s.log, "Close")

	logger.Info("draining nats connection")

	if err := s.conn.Drain(); err != nil {
		logger.Error("draining nats connection", pkglog.Err(err))
		s.conn.Close()
	}

	s.subs = nil
}

func (s *DocumentSubscriber) handleDocumentAnalyzed(msg *nc.Msg) {
	logger := pkglog.WithMethod(s.log, "handleDocumentAnalyzed")

	var event eventuserpb.DocumentAnalyzedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		logger.Error("unmarshalling event", pkglog.Err(err))
		return
	}

	if err := s.handler.OnDocumentAnalyzed(context.Background(), natsdto.DocumentAnalyzedEventFromProto(&event)); err != nil {
		logger.Error("handling event",
			slog.String("documentID", event.GetDocumentId()),
			pkglog.Err(err),
		)
	}
}
