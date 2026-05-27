package model

import (
	sharedmodel "carsharing/shared/model"
	"time"
)

type CarMaintenanceTemplate struct {
	ID   string
	Name string

	KmInterval  *int32
	DayInterval *int32

	IsMandatory bool
	WarnPct     float64
	PullPct     float64
}

type CarMaintenanceRecord struct {
	ID         string
	CarID      string
	TemplateID string

	Status                string
	MileageAtWarningKM    int32
	MileageAtCompletionKM *int32
	CostTenge             *int32
	AssignedTo            *string
	ReceiptImageURLs      []string
	Notes                 *string

	DueBy       *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
}

type CarMaintenanceTemplateFilter struct {
	Pagination *sharedmodel.Pagination
}

type CarMaintenanceRecordFilter struct {
	CarID      *string
	TemplateID *string
	Status     *string

	Pagination *sharedmodel.Pagination
}

type CarMaintenanceTemplateCreate struct {
	Name string

	KmInterval  *int32
	DayInterval *int32

	IsMandatory bool
	WarnPct     float64
	PullPct     float64
}

type CarMaintenanceTemplateUpdate struct {
	Name *string

	KmInterval  *int32
	DayInterval *int32

	IsMandatory *bool
	WarnPct     *float64
	PullPct     *float64
}

type CarMaintenanceRecordComplete struct {
	MileageAtCompletionKM int32
	CostTenge             int32
	ReceiptImageKeys      []string
	Notes                 *string
}

type CarMaintenanceTemplateAssign struct {
	CarID       string
	TemplateID  string
	InitialKM   *int32
	InitialDate *time.Time
}

type CarMaintenanceEvent struct {
	CarID      string
	TemplateID string
	RecordID   string
	EventType  string
	OccurredAt time.Time
}
