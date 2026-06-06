package handler

import (
	"context"
	"errors"
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

type TripWsHandler struct {
	svc TripStreamService
	log *slog.Logger
}

func NewTripWsHandler(svc TripStreamService, logger *slog.Logger) *TripWsHandler {
	return &TripWsHandler{
		svc: svc,
		log: pkglog.WithComponent(logger, "ws.TripHandler"),
	}
}

// LiveFeed godoc
// @Summary      Live trip feed
// @Description  WebSocket stream of in-progress trip data (elapsed time, current cost, distance). Streams until the trip ends or the client disconnects.
// @Tags         trips
// @Security     BearerAuth
// @Param        id  path  string  true  "Trip ID"
// @Produce      json
// @Success      101  {object}  wsdto.TripLiveFeedMessage  "Streamed WebSocket message format"
// @Failure      401  "unauthorized"
// @Failure      500  "internal server error"
// @Router       /ws/trips/{id} [get]
func (h *TripWsHandler) LiveFeed(c *gin.Context) {
	logger := pkglog.WithMethod(h.log, "LiveFeed")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(c.Request.Context()))

	tripID := c.Param("id")

	conn, err := acceptWebSocket(c, nil)
	if err != nil {
		logger.Error("accepting websocket", pkglog.Err(err))
		return
	}
	defer conn.CloseNow()

	ctx, cancel := tokenDeadlineCtx(c)
	defer cancel()
	ctx = conn.CloseRead(ctx)

	streamErr := h.svc.StreamTripLiveFeed(ctx, tripID, func(feed model.TripLiveFeed) error {
		writeCtx, writeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer writeCancel()
		return wsjson.Write(writeCtx, conn, wsdto.TripLiveFeedMessage{
			ElapsedSeconds:     feed.ElapsedSeconds,
			CurrentCostTenge:   feed.CurrentCostTenge,
			DistanceTraveledKM: feed.DistanceTraveledKM,
		})
	})
	switch {
	case streamErr == nil:
		conn.Close(websocket.StatusNormalClosure, "")
	case errors.Is(streamErr, model.ErrForbidden), errors.Is(streamErr, model.ErrUnauthorized):
		conn.Close(websocket.StatusPolicyViolation, streamErr.Error())
	default:
		logger.Error("live feed stream error", pkglog.Err(streamErr))
		conn.Close(websocket.StatusInternalError, "")
	}
}
