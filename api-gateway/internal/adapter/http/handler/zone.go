package handler

import (
	"log/slog"

	"carsharing/api-gateway/internal/adapter/http/dto"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ZoneHandler struct {
	svc ZoneService
	log *slog.Logger
}

func NewZoneHandler(svc ZoneService, log *slog.Logger) *ZoneHandler {
	return &ZoneHandler{
		svc: svc,
		log: pkglog.WithComponent(log, "http.ZoneHandler"),
	}
}

// Create (Zone) godoc
// @Summary      Create zone
// @Description  Creates a geographic zone polygon (operating area, no-drop, parking hub, or surcharge zone).
// @Tags         zones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.ZoneCreateRequest  true  "Zone payload (boundary as GeoJSON polygon string)"
// @Success      200   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /zones [post]
func (h *ZoneHandler) Create(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Create"), utils.MetadataFromCtx(ctx))
	log.Debug("creating zone")

	data, err := dto.FromZoneCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.Create(ctx, data)
	if err != nil {
		log.Warn("creating zone", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	log.Debug("zone created", slog.String("id", id))

	dto.Ok(ctx, gin.H{"id": id})
}

// Get (Zone) godoc
// @Summary      Get zone by ID
// @Tags         zones
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Zone UUID"
// @Success      200  {object}  dto.ZoneResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /zones/{id} [get]
func (h *ZoneHandler) Get(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Get"), utils.MetadataFromCtx(ctx))
	log.Debug("getting zone")

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	zone, err := h.svc.Get(ctx, id)
	if err != nil {
		log.Warn("getting zone", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"zone": dto.ToZoneResponse(zone)})
}

// List (Zone) godoc
// @Summary      List zones
// @Description  Returns zones filtered by type and active status.
// @Tags         zones
// @Produce      json
// @Security     BearerAuth
// @Param        type      query     string   false  "Zone type"  Enums(operating, no_drop, parking_hub, surcharge)
// @Param        isActive  query     boolean  false  "Filter by active flag"
// @Success      200       {object}  dto.ZonesResponse
// @Failure      400       {object}  dto.ErrorResponse
// @Failure      401       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /zones [get]
func (h *ZoneHandler) List(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "List"), utils.MetadataFromCtx(ctx))
	log.Debug("listing zones")

	filter, err := dto.ZoneFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	zones, err := h.svc.List(ctx, filter)
	if err != nil {
		log.Warn("listing zones", pkglog.Err(err))

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
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /zones/{id} [patch]
func (h *ZoneHandler) Update(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Update"), utils.MetadataFromCtx(ctx))
	log.Debug("updating zone")

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

	if err = h.svc.Update(ctx, id, data); err != nil {
		log.Warn("updating zone", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// Delete (Zone) godoc
// @Summary      Delete zone
// @Tags         zones
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Zone UUID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /zones/{id} [delete]
func (h *ZoneHandler) Delete(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "Delete"), utils.MetadataFromCtx(ctx))
	log.Debug("deleting zone")

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if err = h.svc.Delete(ctx, id); err != nil {
		log.Warn("deleting zone", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}
