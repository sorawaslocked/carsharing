package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

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
	ID                      string     `json:"id"`
	CarID                   string     `json:"carId"`
	TemplateID              string     `json:"templateId"`
	Status                  string     `json:"status"`
	OdometerAt              int32      `json:"odometerAt"`
	CompletedKm             *int32     `json:"completedKm,omitempty"`
	CostTenge               *int32     `json:"costTenge,omitempty"`
	AssignedTo              *string    `json:"assignedTo,omitempty"`
	DueBy                   *time.Time `json:"dueBy,omitempty"`
	CompletedAt             *time.Time `json:"completedAt,omitempty"`
	ReceiptImageStorageUrls []string   `json:"receiptImageStorageUrls,omitempty"`
	Notes                   *string    `json:"notes,omitempty"`
	CreatedAt               time.Time  `json:"createdAt"`
}

type CarMaintenanceTemplateCreateRequest struct {
	Name        string  `json:"name"`
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
	CompletedKm             int32    `json:"completedKm"`
	CostTenge               int32    `json:"costTenge"`
	ReceiptImageStorageKeys []string `json:"receiptImageStorageKeys"`
	Notes                   *string  `json:"notes"`
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
		CompletedKm:             req.CompletedKm,
		CostTenge:               req.CostTenge,
		ReceiptImageStorageKeys: req.ReceiptImageStorageKeys,
		Notes:                   req.Notes,
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

	if v := ctx.Param("carID"); v != "" {
		f.CarID = &v
	}
	if v := ctx.Param("templateID"); v != "" {
		f.TemplateID = &v
	}
	if v := ctx.Param("status"); v != "" {
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
		ID:                      m.CarID,
		CarID:                   m.CarID,
		TemplateID:              m.TemplateID,
		Status:                  m.Status,
		OdometerAt:              m.OdometerAt,
		CompletedKm:             m.CompletedKm,
		CostTenge:               m.CostTenge,
		Notes:                   m.Notes,
		AssignedTo:              m.AssignedTo,
		DueBy:                   m.DueBy,
		CompletedAt:             m.CompletedAt,
		ReceiptImageStorageUrls: m.ReceiptImageStorageUrls,
		CreatedAt:               m.CreatedAt,
	}
}
