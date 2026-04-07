package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
)

type ZoneHandler struct {
	svc ZoneService
}

func NewZoneHandler(svc ZoneService) *ZoneHandler {
	return &ZoneHandler{svc: svc}
}

// Create (Zone) godoc
// @Summary      Create zone
// @Description  Creates a geographic zone polygon (operating area, no-drop, parking hub, or surcharge zone).
// @Tags         zones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.ZoneCreateRequest  true  "Zone payload (boundary as GeoJSON polygon string)"
// @Success      200   {object}  map[string]any         "id"
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      409   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /zones [post]
func (h *ZoneHandler) Create(ctx *gin.Context) {
	data, err := dto.FromZoneCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"id": id})
}

// Get (Zone) godoc
// @Summary      Get zone by ID
// @Tags         zones
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Zone UUID"
// @Success      200  {object}  map[string]any  "zone"
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /zones/{id} [get]
func (h *ZoneHandler) Get(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	zone, err := h.svc.Get(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"zone": dto.ToZoneResponse(zone)})
}

// GetAll (Zone) godoc
// @Summary      List zones
// @Description  Returns zones filtered by type and active status.
// @Tags         zones
// @Produce      json
// @Security     BearerAuth
// @Param        type      query     string   false  "Zone type (operating, no_drop, parking_hub, surcharge)"
// @Param        isActive  query     boolean  false  "Filter by active flag"
// @Success      200       {object}  map[string]any  "zones"
// @Failure      400       {object}  map[string]any
// @Failure      500       {object}  map[string]any
// @Router       /zones [get]
func (h *ZoneHandler) GetAll(ctx *gin.Context) {
	filter, err := dto.ZoneFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	zones, err := h.svc.GetAll(ctx, filter)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	zoneResponse := make([]dto.Zone, len(zones))
	for i, zone := range zones {
		zoneResponse[i] = dto.ToZoneResponse(zone)
	}

	dto.Ok(ctx, gin.H{"zones": zoneResponse})
}

// Update (Zone) godoc
// @Summary      Update zone
// @Tags         zones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                 true  "Zone UUID"
// @Param        body  body      dto.ZoneUpdateRequest  true  "Fields to update"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      404   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /zones/{id} [patch]
func (h *ZoneHandler) Update(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromZoneUpdateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	err = h.svc.Update(ctx, id, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}

// Delete (Zone) godoc
// @Summary      Delete zone
// @Tags         zones
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Zone UUID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /zones/{id} [delete]
func (h *ZoneHandler) Delete(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	err = h.svc.Delete(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}
