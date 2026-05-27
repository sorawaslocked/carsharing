package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type MaintenanceRecordStatus string

const (
	MaintenanceRecordStatusPending    MaintenanceRecordStatus = "pending"
	MaintenanceRecordStatusInProgress MaintenanceRecordStatus = "in_progress"
	MaintenanceRecordStatusCompleted  MaintenanceRecordStatus = "completed"
)

var validMaintenanceRecordStatuses = map[MaintenanceRecordStatus]struct{}{
	MaintenanceRecordStatusPending:    {},
	MaintenanceRecordStatusInProgress: {},
	MaintenanceRecordStatusCompleted:  {},
}

func MaintenanceRecordStatusFromString(s string) (MaintenanceRecordStatus, bool) {
	ms := MaintenanceRecordStatus(s)
	if _, ok := validMaintenanceRecordStatuses[ms]; !ok {
		return "", false
	}
	return ms, true
}

func (s MaintenanceRecordStatus) String() string {
	return string(s)
}

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

	MileageAtWarningKM int32
	CompletedKM        *int32
	CostTenge          *int32
	AssignedTo         *string

	DueBy       *time.Time
	CompletedAt *time.Time

	Notes         *string
	ReceiptImages []sharedmodel.Image

	CreatedAt time.Time
	UpdatedAt time.Time
}
type CarMaintenanceEvaluation struct {
	CarID      string
	TemplateID string
	Pct        float64
	Action     string
}

type CarMaintenanceEvent struct {
	CarID      string
	TemplateID string
	RecordID   string
	EventType  string // "warn" or "pull"
	OccurredAt time.Time
}

type CarMaintenanceTemplateFilter struct {
	IsMandatory *bool

	Pagination *sharedmodel.Pagination
}

type CarMaintenanceRecordFilter struct {
	CarID      *string
	TemplateID *string
	Status     *MaintenanceRecordStatus

	Pagination *sharedmodel.Pagination
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
