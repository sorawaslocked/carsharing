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
}

type CarMaintenanceRecord struct {
	ID                      string
	CarID                   string
	TemplateID              string
	Status                  string
	OdometerAt              int32
	CompletedKm             *int32
	CostTenge               *int32
	AssignedTo              *string
	DueBy                   *time.Time
	CompletedAt             *time.Time
	ReceiptImageStorageUrls []string
	Notes                   *string
	CreatedAt               time.Time
}

type CarMaintenanceTemplateFilter struct {
	Pagination *Pagination
}

type CarMaintenanceRecordFilter struct {
	CarID      *string
	TemplateID *string
	Status     *string
	Pagination *Pagination
}

type CarMaintenanceTemplateCreate struct {
	Name        string
	KmInterval  *int32
	DayInterval *int32
	IsMandatory bool
	WarnPct     float64
	PullPct     float64
}

type CarMaintenanceTemplateUpdate struct {
	Name        *string
	KmInterval  *int32
	DayInterval *int32
	IsMandatory *bool
	WarnPct     *float64
	PullPct     *float64
}

type CarMaintenanceRecordComplete struct {
	CompletedKm             int32
	CostTenge               int32
	ReceiptImageStorageKeys []string
	Notes                   *string
}
