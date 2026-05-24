package handler

import (
	"log/slog"

	httpdto "carsharing/api-gateway/internal/adapter/http/dto"
	wsdto "carsharing/api-gateway/internal/adapter/websocket/dto"
	"carsharing/api-gateway/internal/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gin-gonic/gin"
)

type CarWsHandler struct {
	svc          CarStreamService
	carStatusHub *CarStatusHub
	log          *slog.Logger
}

func NewCarWsHandler(svc CarStreamService, carStatusHub *CarStatusHub, logger *slog.Logger) *CarWsHandler {
	return &CarWsHandler{
		svc:          svc,
		carStatusHub: carStatusHub,
		log:          pkglog.WithComponent(logger, "ws.CarHandler"),
	}
}

// Fleet godoc
// @Summary      Live car fleet feed
// @Description  WebSocket stream of slim car objects matching the filter. Accepts the same query params as GET /cars. Streams updates until token expiry or disconnect.
// @Tags         cars
// @Security     BearerAuth
// @Param        brand         query  string   false  "Filter by brand"
// @Param        model         query  string   false  "Filter by model"
// @Param        fuelType      query  string   false  "Filter by fuel type"
// @Param        transmission  query  string   false  "Filter by transmission"
// @Param        bodyType      query  string   false  "Filter by body type"
// @Param        class         query  string   false  "Filter by class"
// @Param        minSeats      query  integer  false  "Minimum seats"
// @Param        latitude      query  number   false  "Center latitude for radius search"
// @Param        longitude     query  number   false  "Center longitude for radius search"
// @Param        radiusM       query  integer  false  "Radius in meters"
// @Param        zoneId        query  string   false  "Filter by zone ID"
// @Param        minFuelLevel  query  number   false  "Minimum fuel level"
// @Param        status        query  string   false  "Filter by status"
// @Produce      json
// @Success      101  {object}  wsdto.CarFleetMessage  "Streamed WebSocket message format"
// @Failure      400  "bad request"
// @Failure      401  "unauthorized"
// @Failure      500  "internal server error"
// @Router       /ws/cars [get]
func (h *CarWsHandler) Fleet(c *gin.Context) {
	logger := pkglog.WithMethod(h.log, "Fleet")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(c.Request.Context()))

	filter, err := httpdto.CarFilterFromCtx(c)
	if err != nil {
		c.Status(400)
		return
	}

	conn, err := acceptWebSocket(c, nil)
	if err != nil {
		logger.Error("accepting websocket", pkglog.Err(err))
		return
	}
	defer conn.CloseNow()

	ctx, cancel := tokenDeadlineCtx(c)
	defer cancel()

	streamErr := h.svc.StreamCarsWithFilter(ctx, filter, func(cars []model.SlimCar) error {
		slim := make([]wsdto.SlimCar, len(cars))
		for i, c := range cars {
			slim[i] = wsdto.SlimCar{
				ID:           c.ID,
				ModelID:      c.ModelID,
				LicensePlate: c.LicensePlate,
				Color:        c.Color,
				Location:     wsdto.SlimCarLocation{Latitude: c.Location.Latitude, Longitude: c.Location.Longitude},
				FuelLevel:    c.FuelLevel,
				Status:       c.Status,
			}
		}
		return wsjson.Write(ctx, conn, wsdto.CarFleetMessage{Cars: slim})
	})
	if streamErr != nil {
		logger.Error("fleet stream error", pkglog.Err(streamErr))
	}

	conn.Close(websocket.StatusNormalClosure, "")
}

// Telemetry godoc
// @Summary      Live car telemetry feed
// @Description  WebSocket stream of telemetry updates (fuel, battery, mileage, location) for a single car. Streams until token expiry or disconnect.
// @Tags         cars
// @Security     BearerAuth
// @Param        id  path  string  true  "Car ID"
// @Produce      json
// @Success      101  {object}  wsdto.CarTelemetryMessage  "Streamed WebSocket message format"
// @Failure      401  "unauthorized"
// @Failure      500  "internal server error"
// @Router       /ws/cars/{id}/telemetry [get]
func (h *CarWsHandler) Telemetry(c *gin.Context) {
	logger := pkglog.WithMethod(h.log, "Telemetry")
	logger = pkglog.WithMetadata(logger, utils.MetadataFromCtx(c.Request.Context()))

	carID := c.Param("id")

	conn, err := acceptWebSocket(c, nil)
	if err != nil {
		logger.Error("accepting websocket", pkglog.Err(err))
		return
	}
	defer conn.CloseNow()

	ctx, cancel := tokenDeadlineCtx(c)
	defer cancel()

	streamErr := h.svc.StreamCarTelemetry(ctx, carID, func(event model.CarTelemetryEvent) error {
		return wsjson.Write(ctx, conn, wsdto.CarTelemetryMessage{
			FuelLevel:    event.FuelLevel,
			BatteryLevel: event.BatteryLevel,
			MileageKM:    event.MileageKM,
			Location:     wsdto.SlimCarLocation{Latitude: event.Location.Latitude, Longitude: event.Location.Longitude},
			RecordedAt:   event.RecordedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	})
	if streamErr != nil {
		logger.Error("telemetry stream error", pkglog.Err(streamErr))
	}

	conn.Close(websocket.StatusNormalClosure, "")
}

// Status godoc
// @Summary      Live car status feed
// @Description  WebSocket stream of status transition events for a single car, delivered via NATS. Streams until token expiry or disconnect.
// @Tags         cars
// @Security     BearerAuth
// @Param        id  path  string  true  "Car ID"
// @Produce      json
// @Success      101  {object}  wsdto.CarStatusMessage  "Streamed WebSocket message format"
// @Failure      401  "unauthorized"
// @Failure      500  "internal server error"
// @Router       /ws/cars/{id}/status [get]
func (h *CarWsHandler) Status(c *gin.Context) {
	logger := pkglog.WithMethod(h.log, "Status")

	carID := c.Param("id")

	conn, err := acceptWebSocket(c, nil)
	if err != nil {
		logger.Error("accepting websocket", pkglog.Err(err))
		return
	}
	defer conn.CloseNow()

	ctx, cancel := tokenDeadlineCtx(c)
	defer cancel()

	ch, unsub := h.carStatusHub.Subscribe(carID)
	defer unsub()

	for {
		select {
		case event := <-ch:
			msg := wsdto.CarStatusMessage{
				CarID:      event.CarID,
				FromStatus: event.FromStatus,
				ToStatus:   event.ToStatus,
			}
			if writeErr := wsjson.Write(ctx, conn, msg); writeErr != nil {
				logger.Error("writing message", pkglog.Err(writeErr))
				conn.Close(websocket.StatusNormalClosure, "")
				return
			}

		case <-ctx.Done():
			conn.Close(websocket.StatusNormalClosure, "")
			return
		}
	}
}
