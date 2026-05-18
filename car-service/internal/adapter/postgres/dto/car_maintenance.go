package dto

import (
	"database/sql"
	"fmt"
	"time"

	"carsharing/car-service/internal/model"
	"github.com/lib/pq"
)

// --- CarMaintenanceTemplate ---

type carMaintenanceTemplateRow struct {
	ID          string
	Name        string
	KmInterval  sql.NullInt32
	DayInterval sql.NullInt32
	IsMandatory bool
	WarnPct     float64
	PullPct     float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r carMaintenanceTemplateRow) toDomain() model.CarMaintenanceTemplate {
	t := model.CarMaintenanceTemplate{
		ID:          r.ID,
		Name:        r.Name,
		IsMandatory: r.IsMandatory,
		WarnPct:     r.WarnPct,
		PullPct:     r.PullPct,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	if r.KmInterval.Valid {
		v := r.KmInterval.Int32
		t.KmInterval = &v
	}
	if r.DayInterval.Valid {
		v := r.DayInterval.Int32
		t.DayInterval = &v
	}

	return t
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

func BuildMaintenanceTemplateWhereClauses(f model.CarMaintenanceTemplateFilter, b *ArgsBuilder) []string {
	var clauses []string

	if f.IsMandatory != nil {
		clauses = append(clauses, fmt.Sprintf("is_mandatory = %s", b.Add(*f.IsMandatory)))
	}

	return clauses
}

func BuildMaintenanceTemplateSetClauses(u model.CarMaintenanceTemplateUpdate, b *ArgsBuilder) []string {
	var clauses []string

	if u.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = %s", b.Add(*u.Name)))
	}
	if u.KmInterval != nil {
		clauses = append(clauses, fmt.Sprintf("km_interval = %s", b.Add(*u.KmInterval)))
	}
	if u.DayInterval != nil {
		clauses = append(clauses, fmt.Sprintf("day_interval = %s", b.Add(*u.DayInterval)))
	}
	if u.IsMandatory != nil {
		clauses = append(clauses, fmt.Sprintf("is_mandatory = %s", b.Add(*u.IsMandatory)))
	}
	if u.WarnPct != nil {
		clauses = append(clauses, fmt.Sprintf("warn_pct = %s", b.Add(*u.WarnPct)))
	}
	if u.PullPct != nil {
		clauses = append(clauses, fmt.Sprintf("pull_pct = %s", b.Add(*u.PullPct)))
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = %s", b.Add(u.UpdatedAt)))

	return clauses
}

// --- CarMaintenanceRecord ---

type carMaintenanceRecordRow struct {
	ID          string
	CarID       string
	TemplateID  string
	Status      string
	OdometerAt  int32
	CompletedKM sql.NullInt32
	CostTenge   sql.NullInt32
	AssignedTo  sql.NullString
	DueBy       sql.NullTime
	CompletedAt sql.NullTime
	Notes       sql.NullString
	ReceiptKeys pq.StringArray
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r carMaintenanceRecordRow) toDomain() model.CarMaintenanceRecord {
	rec := model.CarMaintenanceRecord{
		ID:            r.ID,
		CarID:         r.CarID,
		TemplateID:    r.TemplateID,
		Status:        model.MaintenanceRecordStatus(r.Status),
		OdometerAt:    r.OdometerAt,
		ReceiptImages: ImageKeysToImages([]string(r.ReceiptKeys)),
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}

	if r.CompletedKM.Valid {
		v := r.CompletedKM.Int32
		rec.CompletedKM = &v
	}
	if r.CostTenge.Valid {
		v := r.CostTenge.Int32
		rec.CostTenge = &v
	}
	if r.AssignedTo.Valid {
		rec.AssignedTo = &r.AssignedTo.String
	}
	if r.DueBy.Valid {
		rec.DueBy = &r.DueBy.Time
	}
	if r.CompletedAt.Valid {
		rec.CompletedAt = &r.CompletedAt.Time
	}
	if r.Notes.Valid {
		rec.Notes = &r.Notes.String
	}

	return rec
}

func ScanMaintenanceRecordRow(s scanner) (model.CarMaintenanceRecord, error) {
	var r carMaintenanceRecordRow

	err := s.Scan(
		&r.ID, &r.CarID, &r.TemplateID, &r.Status, &r.OdometerAt,
		&r.CompletedKM, &r.CostTenge, &r.AssignedTo,
		&r.DueBy, &r.CompletedAt, &r.Notes, &r.ReceiptKeys,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.CarMaintenanceRecord{}, err
	}

	return r.toDomain(), nil
}

func BuildMaintenanceRecordWhereClauses(f model.CarMaintenanceRecordFilter, b *ArgsBuilder) []string {
	var clauses []string

	if f.CarID != nil {
		clauses = append(clauses, fmt.Sprintf("car_id = %s", b.Add(*f.CarID)))
	}
	if f.TemplateID != nil {
		clauses = append(clauses, fmt.Sprintf("template_id = %s", b.Add(*f.TemplateID)))
	}
	if f.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = %s", b.Add(string(*f.Status))))
	}

	return clauses
}

func BuildMaintenanceRecordSetClauses(u model.CarMaintenanceRecordUpdate, b *ArgsBuilder) []string {
	var clauses []string

	if u.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = %s", b.Add(string(*u.Status))))
	}
	if u.AssignedTo != nil {
		clauses = append(clauses, fmt.Sprintf("assigned_to = %s", b.Add(*u.AssignedTo)))
	}
	if u.DueBy != nil {
		clauses = append(clauses, fmt.Sprintf("due_by = %s", b.Add(*u.DueBy)))
	}
	if u.CompletedKM != nil {
		clauses = append(clauses, fmt.Sprintf("completed_km = %s", b.Add(*u.CompletedKM)))
	}
	if u.CostTenge != nil {
		clauses = append(clauses, fmt.Sprintf("cost_tenge = %s", b.Add(*u.CostTenge)))
	}
	if u.CompletedAt != nil {
		clauses = append(clauses, fmt.Sprintf("completed_at = %s", b.Add(*u.CompletedAt)))
	}
	if u.Notes != nil {
		clauses = append(clauses, fmt.Sprintf("notes = %s", b.Add(*u.Notes)))
	}
	if u.ReceiptImageKeys != nil {
		clauses = append(clauses, fmt.Sprintf("receipt_image_keys = %s", b.Add(pq.StringArray(u.ReceiptImageKeys))))
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = %s", b.Add(u.UpdatedAt)))

	return clauses
}
