package handler

import (
	"log/slog"

	wsdto "carsharing/api-gateway/internal/adapter/websocket/dto"
	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gin-gonic/gin"
)

type UserWsHandler struct {
	svc DocumentStreamService
	log *slog.Logger
}

func NewUserWsHandler(svc DocumentStreamService, logger *slog.Logger) *UserWsHandler {
	return &UserWsHandler{
		svc: svc,
		log: pkglog.WithComponent(logger, "ws.UserHandler"),
	}
}

// DocumentUpdates godoc
// @Summary      Document analyzed stream
// @Description  WebSocket stream of DocumentAnalyzedEvents from the user service. Accepts optional userId and passed query params to filter events.
// @Tags         users
// @Security     BearerAuth
// @Param        userId  query  string  false  "Filter by user ID"
// @Param        passed  query  bool    false  "Filter by pass/fail result"
// @Produce      json
// @Success      101  {object}  wsdto.DocumentAnalyzedMessage  "Streamed WebSocket message format"
// @Failure      401  "unauthorized"
// @Failure      500  "internal server error"
// @Router       /ws/users/documents [get]
func (h *UserWsHandler) DocumentUpdates(c *gin.Context) {
	logger := pkglog.WithMethod(h.log, "DocumentUpdates")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(c.Request.Context()))

	var userID *string
	if v := c.Query("userID"); v != "" {
		userID = &v
	}

	var passed *bool
	if v := c.Query("passed"); v != "" {
		b := v == "true"
		passed = &b
	}

	conn, err := acceptWebSocket(c, nil)
	if err != nil {
		logger.Error("accepting websocket", pkglog.Err(err))
		return
	}
	defer conn.CloseNow()

	ctx, cancel := tokenDeadlineCtx(c)
	defer cancel()

	streamErr := h.svc.StreamDocumentAnalyzed(ctx, userID, passed, func(event model.DocumentAnalyzedEvent) error {
		defects := make([]wsdto.DocumentDefect, len(event.Defects))
		for i, d := range event.Defects {
			defects[i] = wsdto.DocumentDefect{
				Type:        d.Type,
				Description: d.Description,
			}
		}
		return wsjson.Write(ctx, conn, wsdto.DocumentAnalyzedMessage{
			DocumentID: event.DocumentID,
			UserID:     event.UserID,
			Passed:     event.Passed,
			Defects:    defects,
		})
	})
	if streamErr != nil {
		logger.Error("document stream error", pkglog.Err(streamErr))
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
