package handler

import (
	"log/slog"
	"net/http"

	wsdto "carsharing/api-gateway/internal/adapter/websocket/dto"
	pkglog "carsharing/shared/pkg/log"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gin-gonic/gin"
)

type UserWsHandler struct {
	hub *DocumentHub
	log *slog.Logger
}

func NewUserWsHandler(hub *DocumentHub, logger *slog.Logger) *UserWsHandler {
	return &UserWsHandler{
		hub: hub,
		log: pkglog.WithComponent(logger, "ws.UserHandler"),
	}
}

// DocumentUpdates godoc
// @Summary      Document verification updates
// @Description  WebSocket feed that pushes a single message when the specified document has been analyzed, then closes.
// @Tags         users
// @Security     BearerAuth
// @Param        docID  query  string  true  "Document ID to watch"
// @Produce      json
// @Success      101  {object}  wsdto.DocumentAnalyzedMessage  "WebSocket message (one-shot, then closes)"
// @Failure      400  "bad request"
// @Failure      401  "unauthorized"
// @Failure      500  "internal server error"
// @Router       /ws/users/documents [get]
func (h *UserWsHandler) DocumentUpdates(c *gin.Context) {
	logger := pkglog.WithMethod(h.log, "DocumentUpdates")

	docID := c.Query("docID")
	if docID == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("accepting websocket", pkglog.Err(err))
		return
	}
	defer conn.CloseNow()

	ctx := c.Request.Context()

	ch, unsub := h.hub.Subscribe(docID)
	defer unsub()

	select {
	case event := <-ch:
		defects := make([]wsdto.DocumentDefect, len(event.Defects))
		for i, d := range event.Defects {
			defects[i] = wsdto.DocumentDefect{
				Type:        d.Type,
				Description: d.Description,
			}
		}

		msg := wsdto.DocumentAnalyzedMessage{
			DocumentID: event.DocumentID,
			Passed:     event.Passed,
			Defects:    defects,
		}

		if writeErr := wsjson.Write(ctx, conn, msg); writeErr != nil {
			logger.Error("writing message", pkglog.Err(writeErr))
		}

		conn.Close(websocket.StatusNormalClosure, "")

	case <-ctx.Done():
		conn.Close(websocket.StatusNormalClosure, "")
	}
}
