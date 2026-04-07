package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/adapter/http/dto"
)

type CarMaintenanceHandler struct {
	svc CarMaintenanceService
}

func NewCarMaintenanceHandler(svc CarMaintenanceService) *CarMaintenanceHandler {
	return &CarMaintenanceHandler{svc: svc}
}

// CreateTemplate godoc
// @Summary      Create maintenance template
// @Description  Defines a recurring service interval (e.g. oil change every 8,000 km / 180 days).
// @Tags         car-maintenance
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CarMaintenanceTemplateCreateRequest  true  "Template payload"
// @Success      200   {object}  map[string]any                           "id"
// @Failure      400   {object}  map[string]any
// @Failure      401   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /car-maintenance/template [post]
func (h *CarMaintenanceHandler) CreateTemplate(ctx *gin.Context) {
	data, err := dto.FromCarMaintenanceTemplateCreateRequest(ctx)
	if err != nil {
		dto.MalformedJson(ctx)

		return
	}

	id, err := h.svc.CreateTemplate(ctx, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"id": id})
}

// GetTemplate godoc
// @Summary      Get maintenance template by ID
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Template UUID"
// @Success      200  {object}  map[string]any  "template"
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-maintenance/template/{id} [get]
func (h *CarMaintenanceHandler) GetTemplate(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	template, err := h.svc.GetTemplate(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"template": dto.ToCarMaintenanceTemplateResponse(template)})
}

// GetAllTemplates godoc
// @Summary      List maintenance templates
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Param        limit   query     integer  false  "Pagination limit"
// @Param        offset  query     integer  false  "Pagination offset"
// @Success      200     {object}  map[string]any  "templates"
// @Failure      400     {object}  map[string]any
// @Failure      500     {object}  map[string]any
// @Router       /car-maintenance/template [get]
func (h *CarMaintenanceHandler) GetAllTemplates(ctx *gin.Context) {
	filter, err := dto.CarMaintenanceTemplateFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	templates, err := h.svc.GetAllTemplates(ctx, filter)
	if err != nil {
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
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      404   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /car-maintenance/template/{id} [patch]
func (h *CarMaintenanceHandler) UpdateTemplate(ctx *gin.Context) {
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

	err = h.svc.UpdateTemplate(ctx, id, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}

// DeleteTemplate godoc
// @Summary      Delete maintenance template
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Template UUID"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-maintenance/template/{id} [delete]
func (h *CarMaintenanceHandler) DeleteTemplate(ctx *gin.Context) {
	id, err := dto.IDParam(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	err = h.svc.DeleteTemplate(ctx, id)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}

// GetRecords godoc
// @Summary      List maintenance records
// @Description  Returns work orders for cars. Filter by car, template, or status (pending, in_progress, completed).
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Param        carID       query     string   false  "Filter by car UUID"
// @Param        templateID  query     string   false  "Filter by template UUID"
// @Param        status      query     string   false  "Filter by status (pending, in_progress, completed)"
// @Param        limit       query     integer  false  "Pagination limit"
// @Param        offset      query     integer  false  "Pagination offset"
// @Success      200         {object}  map[string]any  "records"
// @Failure      400         {object}  map[string]any
// @Failure      500         {object}  map[string]any
// @Router       /car-maintenance/records [get]
func (h *CarMaintenanceHandler) GetRecords(ctx *gin.Context) {
	filter, err := dto.CarMaintenanceRecordFilterFromCtx(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	records, err := h.svc.GetRecords(ctx, filter)
	if err != nil {
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
// @Param        id    path      string                                    true  "Maintenance record UUID"
// @Param        body  body      dto.CarMaintenanceRecordCompleteRequest   true  "Completion details"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  map[string]any
// @Failure      404   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /car-maintenance/records/complete/{id} [post]
func (h *CarMaintenanceHandler) CompleteRecord(ctx *gin.Context) {
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

	err = h.svc.CompleteRecord(ctx, id, data)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, nil)
}

// GetReceiptImageUploadUrl godoc
// @Summary      Get pre-signed upload URL for maintenance receipts
// @Tags         car-maintenance
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any  "uploadData"
// @Failure      401  {object}  map[string]any
// @Failure      500  {object}  map[string]any
// @Router       /car-maintenance/records/receipt-image-upload [get]
func (h *CarMaintenanceHandler) GetReceiptImageUploadUrl(ctx *gin.Context) {
	uploadData, err := h.svc.GetReceiptImageUploadData(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
