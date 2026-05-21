package handler

import (
	"context"
	"log/slog"

	pkglog "carsharing/shared/pkg/log"
	natsdto "carsharing/user-service/internal/adapter/nats/dto"

	"github.com/nats-io/nats.go"
	eventuserpb "github.com/sorawaslocked/car-rental-protos/gen/event/user"
	"google.golang.org/protobuf/proto"
)

const subjectDocumentAnalyzed = "document.analyzed"

type DocumentHandler struct {
	log     *slog.Logger
	conn    *nats.Conn
	service DocumentService
}

func NewDocumentHandler(log *slog.Logger, conn *nats.Conn, service DocumentService) *DocumentHandler {
	return &DocumentHandler{
		log:     pkglog.WithComponent(log, "adapter.nats.DocumentHandler"),
		conn:    conn,
		service: service,
	}
}

func (h *DocumentHandler) Subscribe() error {
	_, err := h.conn.Subscribe(subjectDocumentAnalyzed, h.handleDocumentAnalyzed)

	return err
}

func (h *DocumentHandler) handleDocumentAnalyzed(msg *nats.Msg) {
	ctx := context.Background()
	log := pkglog.WithMethod(h.log, "handleDocumentAnalyzed")

	var event eventuserpb.DocumentAnalyzedEvent
	if err := proto.Unmarshal(msg.Data, &event); err != nil {
		log.Error("unmarshalling document analyzed event", pkglog.Err(err))

		return
	}

	if err := h.service.HandleDocumentAnalyzed(ctx, natsdto.DocumentAnalyzedEventFromProto(&event)); err != nil {
		log.Error("handling document analyzed event",
			slog.String("documentID", event.GetDocumentId()),
			pkglog.Err(err),
		)
	}
}
