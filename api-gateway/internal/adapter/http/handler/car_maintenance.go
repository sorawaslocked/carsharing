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

func (h *CarMaintenanceHandler) GetReceiptImageUploadUrl(ctx *gin.Context) {
	uploadData, err := h.svc.GetReceiptImageUploadData(ctx)
	if err != nil {
		dto.FromError(ctx, err)

		return
	}

	dto.Ok(ctx, gin.H{"uploadData": dto.ToImageUploadDataResponse(uploadData)})
}
