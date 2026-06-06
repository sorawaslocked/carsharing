package handler

import (
	"log/slog"

	"carsharing/api-gateway/internal/adapter/http/dto"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CarMaintenanceHandler struct {
	svc CarMaintenanceService
	log *slog.Logger
}

func NewCarMaintenanceHandler(svc CarMaintenanceService, log *slog.Logger) *CarMaintenanceHandler {
	return &CarMaintenanceHandler{
		svc: svc,
		log: pkglog.WithComponent(log, "http.CarMaintenanceHandler"),
	}
}

// CreateTemplate godoc
// @Summary      Create maintenance template
// @Description  Defines a recurring service interval (e.g. oil change every 8,000 km / 180 days). kmInterval min 100 km; dayInterval min 1 day; warnPct and pullPct are proportions 0.0–1.0 (not percent), pullPct must be ≥ warnPct.
// @Tags         car-maintenance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CarMaintenanceTemplateCreateRequest  true  "Template payload"
// @Success      200   {object}  dto.IDResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /car-maintenance/template [post]
func (h *CarMaintenanceHandler) CreateTemplate(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CreateTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("creating maintenance template")

	data, err := dto.FromCarMaintenanceTemplateCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.CreateTemplate(ctx, data)
	if err != nil {
		log.Warn("creating maintenance template", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	log.Debug("maintenance template created", slog.String("id", id))

	dto.Ok(ctx, gin.H{"id": id})
}

// GetTemplate godoc
// @Summary      Get maintenance template by ID
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Template UUID"
// @Success      200  {object}  dto.CarMaintenanceTemplateItemResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-maintenance/template/{id} [get]
func (h *CarMaintenanceHandler) GetTemplate(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("getting maintenance template")

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	template, err := h.svc.GetTemplate(ctx, id)
	if err != nil {
		log.Warn("getting maintenance template", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"template": dto.ToCarMaintenanceTemplateResponse(template)})
}

// ListTemplates godoc
// @Summary      List maintenance templates
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Param        limit   query     integer  false  "Pagination limit"
// @Param        offset  query     integer  false  "Pagination offset"
// @Success      200     {object}  dto.CarMaintenanceTemplatesResponse
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      401     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Router       /car-maintenance/template [get]
func (h *CarMaintenanceHandler) ListTemplates(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "ListTemplates"), utils.MetadataFromCtx(ctx))
	log.Debug("listing maintenance templates")

	filter, err := dto.CarMaintenanceTemplateFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	templates, err := h.svc.ListTemplates(ctx, filter)
	if err != nil {
		log.Warn("listing maintenance templates", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	templateResponse := make([]dto.CarMaintenanceTemplate, len(templates))
	for i, template := range templates {
		templateResponse[i] = dto.ToCarMaintenanceTemplateResponse(template)
	}

	dto.Ok(ctx, gin.H{"templates": templateResponse})
}

// UpdateTemplate godoc
// @Summary      Update maintenance template
// @Tags         car-maintenance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                                   true  "Template UUID"
// @Param        body  body      dto.CarMaintenanceTemplateUpdateRequest  true  "Fields to update"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-maintenance/template/{id} [patch]
func (h *CarMaintenanceHandler) UpdateTemplate(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "UpdateTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("updating maintenance template")

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromCarMaintenanceTemplateUpdateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.UpdateTemplate(ctx, id, data); err != nil {
		log.Warn("updating maintenance template", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// DeleteTemplate godoc
// @Summary      Delete maintenance template
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Template UUID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-maintenance/template/{id} [delete]
func (h *CarMaintenanceHandler) DeleteTemplate(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "DeleteTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("deleting maintenance template")

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	if err = h.svc.DeleteTemplate(ctx, id); err != nil {
		log.Warn("deleting maintenance template", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// AssignTemplate godoc
// @Summary      Assign a maintenance template to a car
// @Description  Seeds the service state baseline for a (car, template) pair. Optionally accepts an initial KM and date to reflect prior service history.
// @Tags         car-maintenance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CarMaintenanceTemplateAssignRequest  true  "Assignment payload"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-maintenance/template/assign [post]
func (h *CarMaintenanceHandler) AssignTemplate(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "AssignTemplate"), utils.MetadataFromCtx(ctx))
	log.Debug("assigning maintenance template")

	data, err := dto.FromCarMaintenanceTemplateAssignRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.AssignTemplate(ctx, data); err != nil {
		log.Warn("assigning maintenance template", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// GetRecords godoc
// @Summary      List maintenance records
// @Description  Returns work orders for cars. Filter by car, template, or status (pending, in_progress, completed).
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Param        carID       query     string   false  "Filter by car UUID"
// @Param        templateID  query     string   false  "Filter by template UUID"
// @Param        status      query     string   false  "Filter by status"  Enums(pending, in_progress, completed)
// @Param        limit       query     integer  false  "Pagination limit"
// @Param        offset      query     integer  false  "Pagination offset"
// @Success      200         {object}  dto.CarMaintenanceRecordsResponse
// @Failure      400         {object}  dto.ErrorResponse
// @Failure      401         {object}  dto.ErrorResponse
// @Failure      500         {object}  dto.ErrorResponse
// @Router       /car-maintenance/records [get]
func (h *CarMaintenanceHandler) GetRecords(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetRecords"), utils.MetadataFromCtx(ctx))
	log.Debug("listing maintenance records")

	filter, err := dto.CarMaintenanceRecordFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	records, err := h.svc.ListRecords(ctx, filter)
	if err != nil {
		log.Warn("listing maintenance records", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	recordResponse := make([]dto.CarMaintenanceRecord, len(records))
	for i, record := range records {
		recordResponse[i] = dto.ToCarMaintenanceRecordResponse(record)
	}

	dto.Ok(ctx, gin.H{"records": recordResponse})
}

// CompleteRecord godoc
// @Summary      Complete a maintenance work order
// @Description  Marks a maintenance record as completed, resets the service interval clock, and transitions the car back to available.
// @Tags         car-maintenance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                                   true  "Maintenance record UUID"
// @Param        body  body      dto.CarMaintenanceRecordCompleteRequest  true  "Completion details"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-maintenance/records/complete/{id} [post]
func (h *CarMaintenanceHandler) CompleteRecord(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "CompleteRecord"), utils.MetadataFromCtx(ctx))
	log.Debug("completing maintenance record")

	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	data, err := dto.FromCarMaintenanceRecordCompleteRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	if err = h.svc.CompleteRecord(ctx, id, data); err != nil {
		log.Warn("completing maintenance record", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.NoContent(ctx)
}

// GetReceiptImageUploadUrl godoc
// @Summary      Get pre-signed upload URL for maintenance receipts
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.ImageUploadResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /car-maintenance/records/receipt-image-upload [get]
func (h *CarMaintenanceHandler) GetReceiptImageUploadUrl(ctx *gin.Context) {
	log := pkglog.WithMetadata(pkglog.WithMethod(h.log, "GetReceiptImageUploadUrl"), utils.MetadataFromCtx(ctx))
	log.Debug("getting maintenance receipt image upload url")

	uploadData, err := h.svc.GetReceiptImageUploadData(ctx)
	if err != nil {
		log.Warn("getting maintenance receipt image upload data", pkglog.Err(err))

		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
