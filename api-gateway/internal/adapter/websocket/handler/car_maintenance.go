package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	wsdto "carsharing/api-gateway/internal/adapter/websocket/dto"
	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gin-gonic/gin"
)

type CarMaintenanceWsHandler struct {
	svc CarMaintenanceStreamService
	log *slog.Logger
}

func NewCarMaintenanceWsHandler(svc CarMaintenanceStreamService, logger *slog.Logger) *CarMaintenanceWsHandler {
	return &CarMaintenanceWsHandler{
		svc: svc,
		log: pkglog.WithComponent(logger, "ws.CarMaintenanceHandler"),
	}
}

// MaintenanceEvents godoc
// @Summary      Live maintenance event feed
// @Description  WebSocket stream of maintenance events ("warn" or "pull") emitted by the background evaluator. Streams until token expiry or disconnect.
// @Tags         car-maintenance
// @Security     BearerAuth
// @Produce      json
// @Success      101  {object}  wsdto.CarMaintenanceEventMessage  "Streamed WebSocket message format"
// @Failure      401  "unauthorized"
// @Failure      500  "internal server error"
// @Router       /ws/car-maintenance/events [get]
func (h *CarMaintenanceWsHandler) MaintenanceEvents(c *gin.Context) {
	logger := pkglog.WithMethod(h.log, "MaintenanceEvents")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(c.Request.Context()))

	conn, err := acceptWebSocket(c, nil)
	if err != nil {
		logger.Error("accepting websocket", pkglog.Err(err))
		return
	}
	defer conn.CloseNow()

	connID := fmt.Sprintf("%p", conn)
	logger.Info("maintenance events websocket opened", slog.String("connID", connID))
	defer logger.Info("maintenance events websocket closed", slog.String("connID", connID))

	ctx, cancel := tokenDeadlineCtx(c)
	defer cancel()

	streamErr := h.svc.StreamMaintenanceEvents(ctx, func(event model.CarMaintenanceEvent) error {
		logger.Info("maintenance events writing event",
			slog.String("connID", connID),
			slog.String("eventType", event.EventType),
		)
		writeCtx, writeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer writeCancel()
		return wsjson.Write(writeCtx, conn, wsdto.CarMaintenanceEventMessage{
			CarID:      event.CarID,
			TemplateID: event.TemplateID,
			RecordID:   event.RecordID,
			EventType:  event.EventType,
			OccurredAt: event.OccurredAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	})
	switch {
	case streamErr == nil:
		conn.Close(websocket.StatusNormalClosure, "")
	case errors.Is(streamErr, model.ErrForbidden), errors.Is(streamErr, model.ErrUnauthorized):
		conn.Close(websocket.StatusPolicyViolation, streamErr.Error())
	default:
		logger.Error("maintenance events stream error", pkglog.Err(streamErr), slog.String("connID", connID))
		conn.Close(websocket.StatusInternalError, "")
	}
}
