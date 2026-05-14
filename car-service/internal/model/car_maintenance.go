package model

import "time"

type CarMaintenanceTemplate struct {
	ID          string
	Name        string
	KmInterval  *int32
	DayInterval *int32
	IsMandatory bool
	WarnPct     float64
	PullPct     float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CarServiceState struct {
	CarID       string
	TemplateID  string
	LastKM      int32
	LastDate    *time.Time
	NextDueKM   *int32
	NextDueDate *time.Time
}

type CarMaintenanceRecord struct {
	ID         string
	CarID      string
	TemplateID string
	Status     MaintenanceRecordStatus

	OdometerAt  int32
	CompletedKM *int32
	CostTenge   *int32
	AssignedTo  *string

	DueBy       *time.Time
	CompletedAt *time.Time

	Notes            *string
	ReceiptImageKeys []string

	CreatedAt time.Time
	UpdatedAt time.Time
}
type CarMaintenanceEvaluation struct {
	CarID      string
	TemplateID string
	Pct        float64
	Action     string
}

type CarMaintenanceTemplateFilter struct {
	IsMandatory *bool

	Pagination
}

type CarMaintenanceRecordFilter struct {
	CarID      *string
	TemplateID *string
	Status     *MaintenanceRecordStatus

	Pagination
}

type CarServiceStateFilter struct {
	CarID      *string
	TemplateID *string
}

type CarMaintenanceTemplateUpdate struct {
	Name        *string
	KmInterval  *int32
	DayInterval *int32
	IsMandatory *bool
	WarnPct     *float64
	PullPct     *float64
	UpdatedAt   time.Time
}

type CarMaintenanceRecordUpdate struct {
	Status           *MaintenanceRecordStatus
	AssignedTo       *string
	DueBy            *time.Time
	CompletedKM      *int32
	CostTenge        *int32
	CompletedAt      *time.Time
	Notes            *string
	ReceiptImageKeys []string
	UpdatedAt        time.Time
}

type CarServiceStateUpdate struct {
	LastKM      int32
	LastDate    time.Time
	NextDueKM   *int32
	NextDueDate *time.Time
}

type CarMaintenanceTemplateFilterInput struct {
	IsMandatory *bool `validate:"omitempty"`

	PaginationInput
}

type CarMaintenanceRecordFilterInput struct {
	CarID      *string `validate:"omitempty,uuid"`
	TemplateID *string `validate:"omitempty,uuid"`
	Status     *string `validate:"omitempty,maintenancerecordstatus"`

	PaginationInput
}

type CarMaintenanceTemplateCreateInput struct {
	Name        string  `validate:"required,min=1,max=100"`
	KmInterval  *int32  `validate:"omitempty,min=100"`
	DayInterval *int32  `validate:"omitempty,min=1"`
	IsMandatory bool    `validate:"omitempty"`
	WarnPct     float64 `validate:"min=0,max=1"`
	PullPct     float64 `validate:"min=0,max=1,gtefield=WarnPct"`
}

type CarMaintenanceTemplateUpdateInput struct {
	Name        *string  `validate:"omitempty,min=1,max=100"`
	KmInterval  *int32   `validate:"omitempty,min=100"`
	DayInterval *int32   `validate:"omitempty,min=1"`
	IsMandatory *bool    `validate:"omitempty"`
	WarnPct     *float64 `validate:"omitempty,min=0,max=1"`
	PullPct     *float64 `validate:"omitempty,min=0,max=1"`
}

type CarMaintenanceRecordCompleteInput struct {
	CompletedKM      int32    `validate:"required,min=0"`
	CostTenge        int32    `validate:"min=0"`
	Notes            *string  `validate:"omitempty,min=1,max=1000"`
	ReceiptImageKeys []string `validate:"omitempty,max=10,dive,min=1"`
}
