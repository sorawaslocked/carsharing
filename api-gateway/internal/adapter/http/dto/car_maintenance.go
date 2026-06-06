package dto

import (
	"time"

	"carsharing/api-gateway/internal/model"
	"github.com/gin-gonic/gin"
)

type CarMaintenanceTemplateItemResponse struct {
	Template CarMaintenanceTemplate `json:"template"`
}

type CarMaintenanceTemplatesResponse struct {
	Templates []CarMaintenanceTemplate `json:"templates"`
}

type CarMaintenanceRecordsResponse struct {
	Records []CarMaintenanceRecord `json:"records"`
}

type CarMaintenanceTemplate struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	KmInterval  *int32  `json:"kmInterval,omitempty"`
	DayInterval *int32  `json:"dayInterval,omitempty"`
	IsMandatory bool    `json:"isMandatory"`
	WarnPct     float64 `json:"warnPct"`
	PullPct     float64 `json:"pullPct"`
}

type CarMaintenanceRecord struct {
	ID                 string     `json:"id"`
	CarID              string     `json:"carID"`
	TemplateID         string     `json:"templateID"`
	Status             string     `json:"status" validate:"oneof=pending in_progress completed"`
	MileageAtWarningKM int32      `json:"mileageAtWarningKm"`
	CompletedKm        *int32     `json:"completedKm,omitempty"`
	CostTenge          *int32     `json:"costTenge,omitempty"`
	AssignedTo         *string    `json:"assignedTo,omitempty"`
	DueBy              *time.Time `json:"dueBy,omitempty"`
	CompletedAt        *time.Time `json:"completedAt,omitempty"`
	ReceiptImages      []Image    `json:"receiptImages,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
}

type CarMaintenanceTemplateAssignRequest struct {
	CarID       string     `json:"carID" binding:"required"`
	TemplateID  string     `json:"templateID" binding:"required"`
	InitialKM   *int32     `json:"initialKM"`
	InitialDate *time.Time `json:"initialDate"`
}

type CarMaintenanceTemplateCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	KmInterval  *int32  `json:"kmInterval"`
	DayInterval *int32  `json:"dayInterval"`
	IsMandatory bool    `json:"isMandatory"`
	WarnPct     float64 `json:"warnPct"`
	PullPct     float64 `json:"pullPct"`
}

type CarMaintenanceTemplateUpdateRequest struct {
	Name        *string  `json:"name"`
	KmInterval  *int32   `json:"kmInterval"`
	DayInterval *int32   `json:"dayInterval"`
	IsMandatory *bool    `json:"isMandatory"`
	WarnPct     *float64 `json:"warnPct"`
	PullPct     *float64 `json:"pullPct"`
}

type CarMaintenanceRecordCompleteRequest struct {
	CompletedKm             int32    `json:"completedKm" validate:"min=0"`
	CostTenge               int32    `json:"costTenge" validate:"min=0"`
	ReceiptImageStorageKeys []string `json:"receiptImageStorageKeys"`
	Notes                   *string  `json:"notes"`
}

func FromCarMaintenanceTemplateAssignRequest(ctx *gin.Context) (model.CarMaintenanceTemplateAssign, error) {
	var req CarMaintenanceTemplateAssignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarMaintenanceTemplateAssign{}, err
	}

	return model.CarMaintenanceTemplateAssign{
		CarID:       req.CarID,
		TemplateID:  req.TemplateID,
		InitialKM:   req.InitialKM,
		InitialDate: req.InitialDate,
	}, nil
}

func FromCarMaintenanceTemplateCreateRequest(ctx *gin.Context) (model.CarMaintenanceTemplateCreate, error) {
	var req CarMaintenanceTemplateCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarMaintenanceTemplateCreate{}, err
	}

	return model.CarMaintenanceTemplateCreate{
		Name:        req.Name,
		KmInterval:  req.KmInterval,
		DayInterval: req.DayInterval,
		IsMandatory: req.IsMandatory,
		WarnPct:     req.WarnPct,
		PullPct:     req.PullPct,
	}, nil
}

func FromCarMaintenanceTemplateUpdateRequest(ctx *gin.Context) (model.CarMaintenanceTemplateUpdate, error) {
	var req CarMaintenanceTemplateUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarMaintenanceTemplateUpdate{}, err
	}

	return model.CarMaintenanceTemplateUpdate{
		Name:        req.Name,
		KmInterval:  req.KmInterval,
		DayInterval: req.DayInterval,
		IsMandatory: req.IsMandatory,
		WarnPct:     req.WarnPct,
		PullPct:     req.PullPct,
	}, nil
}

func FromCarMaintenanceRecordCompleteRequest(ctx *gin.Context) (model.CarMaintenanceRecordComplete, error) {
	var req CarMaintenanceRecordCompleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return model.CarMaintenanceRecordComplete{}, err
	}

	return model.CarMaintenanceRecordComplete{
		MileageAtCompletionKM: req.CompletedKm,
		CostTenge:             req.CostTenge,
		ReceiptImageKeys:      req.ReceiptImageStorageKeys,
		Notes:                 req.Notes,
	}, nil
}

func CarMaintenanceTemplateFilterFromCtx(ctx *gin.Context) (model.CarMaintenanceTemplateFilter, error) {
	f := model.CarMaintenanceTemplateFilter{}

	p, err := pagination(ctx)
	if err != nil {
		return f, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func CarMaintenanceRecordFilterFromCtx(ctx *gin.Context) (model.CarMaintenanceRecordFilter, error) {
	f := model.CarMaintenanceRecordFilter{}

	if v := ctx.Query("carID"); v != "" {
		f.CarID = &v
	}
	if v := ctx.Query("templateID"); v != "" {
		f.TemplateID = &v
	}
	if v := ctx.Query("status"); v != "" {
		f.Status = &v
	}

	p, err := pagination(ctx)
	if err != nil {
		return f, model.ErrInvalidQueryParam
	}

	f.Pagination = p

	return f, nil
}

func ToCarMaintenanceTemplateResponse(m model.CarMaintenanceTemplate) CarMaintenanceTemplate {
	return CarMaintenanceTemplate{
		ID:          m.ID,
		Name:        m.Name,
		KmInterval:  m.KmInterval,
		DayInterval: m.DayInterval,
		IsMandatory: m.IsMandatory,
		WarnPct:     m.WarnPct,
		PullPct:     m.PullPct,
	}
}

func ToCarMaintenanceRecordResponse(m model.CarMaintenanceRecord) CarMaintenanceRecord {
	return CarMaintenanceRecord{
		ID:                 m.ID,
		CarID:              m.CarID,
		TemplateID:         m.TemplateID,
		Status:             m.Status,
		MileageAtWarningKM: m.MileageAtWarningKM,
		CompletedKm:        m.MileageAtCompletionKM,
		CostTenge:          m.CostTenge,
		Notes:              m.Notes,
		AssignedTo:         m.AssignedTo,
		DueBy:              m.DueBy,
		CompletedAt:        m.CompletedAt,
		ReceiptImages:      toImages(m.ReceiptImages),
		CreatedAt:          m.CreatedAt,
	}
}
