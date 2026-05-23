package subscriber

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"
	natsdto "carsharing/user-service/internal/adapter/nats/dto"

	eventuserpb "carsharing/protos/gen/event/user"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const subjectDocumentAnalyzed = "document.analyzed"

type DocumentSubscriber struct {
	log     *slog.Logger
	conn    *nats.Conn
	service DocumentService
}

func NewDocumentSubscriber(log *slog.Logger, conn *nats.Conn, service DocumentService) *DocumentSubscriber {
	return &DocumentSubscriber{
		log:     pkglog.WithComponent(log, "adapter.nats.subscriber.DocumentSubscriber"),
		conn:    conn,
		service: service,
	}
}

func (s *DocumentSubscriber) Subscribe() error {
	_, err := s.conn.Subscribe(subjectDocumentAnalyzed, s.handleDocumentAnalyzed)

	return err
}

func (s *DocumentSubscriber) handleDocumentAnalyzed(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(s.log, "handleDocumentAnalyzed")

	var event eventuserpb.DocumentAnalyzedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		log.Error("unmarshalling document analyzed event", pkglog.Err(err))

		return
	}

	if err := s.service.HandleDocumentAnalyzed(ctx, natsdto.DocumentAnalyzedEventFromProto(&event)); err != nil {
		log.Error("handling document analyzed event",
			slog.String("documentID", event.GetDocumentId()),
			pkglog.Err(err),
		)
	}
}
