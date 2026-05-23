package validation

import sharedvalidation "carsharing/shared/validation"

type CarMaintenanceTemplateFilter struct {
	IsMandatory *bool `validate:"omitempty"`
	Pagination  *sharedvalidation.Pagination
}

type CarMaintenanceRecordFilter struct {
	CarID      *string `validate:"omitempty,uuid"`
	TemplateID *string `validate:"omitempty,uuid"`
	Status     *string `validate:"omitempty,maintenancerecordstatus"`
	Pagination *sharedvalidation.Pagination
}

type CarMaintenanceTemplateCreate struct {
	Name        string  `validate:"required,min=1,max=100"`
	KmInterval  *int32  `validate:"omitempty,min=100"`
	DayInterval *int32  `validate:"omitempty,min=1"`
	IsMandatory bool    `validate:"omitempty"`
	WarnPct     float64 `validate:"min=0,max=1"`
	PullPct     float64 `validate:"min=0,max=1,gtefield=WarnPct"`
}

type CarMaintenanceTemplateUpdate struct {
	Name        *string  `validate:"omitempty,min=1,max=100"`
	KmInterval  *int32   `validate:"omitempty,min=100"`
	DayInterval *int32   `validate:"omitempty,min=1"`
	IsMandatory *bool    `validate:"omitempty"`
	WarnPct     *float64 `validate:"omitempty,min=0,max=1"`
	PullPct     *float64 `validate:"omitempty,min=0,max=1"`
}

type CarMaintenanceRecordComplete struct {
	CompletedKM      int32    `validate:"required,min=0"`
	CostTenge        int32    `validate:"min=0"`
	Notes            *string  `validate:"omitempty,min=1,max=1000"`
	ReceiptImageKeys []string `validate:"omitempty,max=10,dive,min=1"`
}
