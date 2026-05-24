package dto

import (
	"fmt"
	"time"

	"carsharing/car-service/internal/model"
)

// --- CarMaintenanceTemplate ---

type carMaintenanceTemplateRow struct {
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

func (r carMaintenanceTemplateRow) toDomain() model.CarMaintenanceTemplate {
	return model.CarMaintenanceTemplate{
		ID:          r.ID,
		Name:        r.Name,
		KmInterval:  r.KmInterval,
		DayInterval: r.DayInterval,
		IsMandatory: r.IsMandatory,
		WarnPct:     r.WarnPct,
		PullPct:     r.PullPct,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func ScanMaintenanceTemplateRow(s scanner) (model.CarMaintenanceTemplate, error) {
	var r carMaintenanceTemplateRow

	err := s.Scan(
		&r.ID, &r.Name, &r.KmInterval, &r.DayInterval,
		&r.IsMandatory, &r.WarnPct, &r.PullPct, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.CarMaintenanceTemplate{}, err
	}

	return r.toDomain(), nil
}

func WhereClausesFromMaintenanceTemplateFilter(f model.CarMaintenanceTemplateFilter, args []any, n int) ([]string, []any, int) {
	var clauses []string

	if f.IsMandatory != nil {
		n++
		args = append(args, *f.IsMandatory)
		clauses = append(clauses, fmt.Sprintf("is_mandatory = $%d", n))
	}

	return clauses, args, n
}

func SetClausesFromMaintenanceTemplateUpdate(update model.CarMaintenanceTemplateUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	n := 0

	if update.Name != nil {
		n++
		args = append(args, *update.Name)
		clauses = append(clauses, fmt.Sprintf("name = $%d", n))
	}
	if update.KmInterval != nil {
		n++
		args = append(args, *update.KmInterval)
		clauses = append(clauses, fmt.Sprintf("km_interval = $%d", n))
	}
	if update.DayInterval != nil {
		n++
		args = append(args, *update.DayInterval)
		clauses = append(clauses, fmt.Sprintf("day_interval = $%d", n))
	}
	if update.IsMandatory != nil {
		n++
		args = append(args, *update.IsMandatory)
		clauses = append(clauses, fmt.Sprintf("is_mandatory = $%d", n))
	}
	if update.WarnPct != nil {
		n++
		args = append(args, *update.WarnPct)
		clauses = append(clauses, fmt.Sprintf("warn_pct = $%d", n))
	}
	if update.PullPct != nil {
		n++
		args = append(args, *update.PullPct)
		clauses = append(clauses, fmt.Sprintf("pull_pct = $%d", n))
	}

	n++
	args = append(args, update.UpdatedAt)
	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", n))

	return clauses, args, n
}

// --- CarMaintenanceRecord ---

type carMaintenanceRecordRow struct {
	ID                 string
	CarID              string
	TemplateID         string
	Status             string
	MileageAtWarningKM int32
	CompletedKM        *int32
	CostTenge          *int32
	AssignedTo         *string
	DueBy              *time.Time
	CompletedAt        *time.Time
	Notes              *string
	ReceiptKeys        []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (r carMaintenanceRecordRow) toDomain() model.CarMaintenanceRecord {
	return model.CarMaintenanceRecord{
		ID:                 r.ID,
		CarID:              r.CarID,
		TemplateID:         r.TemplateID,
		Status:             model.MaintenanceRecordStatus(r.Status),
		MileageAtWarningKM: r.MileageAtWarningKM,
		CompletedKM:        r.CompletedKM,
		CostTenge:          r.CostTenge,
		AssignedTo:         r.AssignedTo,
		DueBy:              r.DueBy,
		CompletedAt:        r.CompletedAt,
		Notes:              r.Notes,
		ReceiptImages:      ImageKeysToImages(r.ReceiptKeys),
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}

func ScanMaintenanceRecordRow(s scanner) (model.CarMaintenanceRecord, error) {
	var r carMaintenanceRecordRow

	err := s.Scan(
		&r.ID, &r.CarID, &r.TemplateID, &r.Status, &r.MileageAtWarningKM,
		&r.CompletedKM, &r.CostTenge, &r.AssignedTo,
		&r.DueBy, &r.CompletedAt, &r.Notes, &r.ReceiptKeys,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.CarMaintenanceRecord{}, err
	}

	return r.toDomain(), nil
}

func WhereClausesFromMaintenanceRecordFilter(f model.CarMaintenanceRecordFilter, args []any, n int) ([]string, []any, int) {
	var clauses []string

	if f.CarID != nil {
		n++
		args = append(args, *f.CarID)
		clauses = append(clauses, fmt.Sprintf("car_id = $%d", n))
	}
	if f.TemplateID != nil {
		n++
		args = append(args, *f.TemplateID)
		clauses = append(clauses, fmt.Sprintf("template_id = $%d", n))
	}
	if f.Status != nil {
		n++
		args = append(args, string(*f.Status))
		clauses = append(clauses, fmt.Sprintf("status = $%d", n))
	}

	return clauses, args, n
}

func SetClausesFromMaintenanceRecordUpdate(update model.CarMaintenanceRecordUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	n := 0

	if update.Status != nil {
		n++
		args = append(args, string(*update.Status))
		clauses = append(clauses, fmt.Sprintf("status = $%d", n))
	}
	if update.AssignedTo != nil {
		n++
		args = append(args, *update.AssignedTo)
		clauses = append(clauses, fmt.Sprintf("assigned_to = $%d", n))
	}
	if update.DueBy != nil {
		n++
		args = append(args, *update.DueBy)
		clauses = append(clauses, fmt.Sprintf("due_by = $%d", n))
	}
	if update.CompletedKM != nil {
		n++
		args = append(args, *update.CompletedKM)
		clauses = append(clauses, fmt.Sprintf("completed_km = $%d", n))
	}
	if update.CostTenge != nil {
		n++
		args = append(args, *update.CostTenge)
		clauses = append(clauses, fmt.Sprintf("cost_tenge = $%d", n))
	}
	if update.CompletedAt != nil {
		n++
		args = append(args, *update.CompletedAt)
		clauses = append(clauses, fmt.Sprintf("completed_at = $%d", n))
	}
	if update.Notes != nil {
		n++
		args = append(args, *update.Notes)
		clauses = append(clauses, fmt.Sprintf("notes = $%d", n))
	}
	if update.ReceiptImageKeys != nil {
		n++
		args = append(args, update.ReceiptImageKeys)
		clauses = append(clauses, fmt.Sprintf("receipt_image_keys = $%d", n))
	}

	n++
	args = append(args, update.UpdatedAt)
	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", n))

	return clauses, args, n
}
