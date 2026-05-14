package dto

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
)

type carInsuranceRow struct {
	ID        string
	CarID     string
	Type      string
	Provider  string
	PolicyNum string
	StartsAt  time.Time
	ExpiresAt time.Time
	CostTenge int32
	Status    string
	Notes     sql.NullString
	ImageKeys pq.StringArray
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r carInsuranceRow) toDomain() model.CarInsurance {
	ins := model.CarInsurance{
		ID:        r.ID,
		CarID:     r.CarID,
		Type:      model.InsuranceType(r.Type),
		Provider:  r.Provider,
		PolicyNum: r.PolicyNum,
		StartsAt:  r.StartsAt,
		ExpiresAt: r.ExpiresAt,
		CostTenge: r.CostTenge,
		Status:    model.InsuranceStatus(r.Status),
		ImageKeys: []string(r.ImageKeys),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	if r.Notes.Valid {
		ins.Notes = &r.Notes.String
	}

	return ins
}

func ScanCarInsuranceRow(s scanner) (model.CarInsurance, error) {
	var r carInsuranceRow

	err := s.Scan(
		&r.ID, &r.CarID, &r.Type, &r.Provider, &r.PolicyNum,
		&r.StartsAt, &r.ExpiresAt, &r.CostTenge, &r.Status,
		&r.Notes, &r.ImageKeys, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.CarInsurance{}, err
	}

	return r.toDomain(), nil
}

func BuildCarInsuranceWhereClauses(f model.CarInsuranceFilter, b *ArgsBuilder) []string {
	var clauses []string

	if f.CarID != nil {
		clauses = append(clauses, fmt.Sprintf("car_id = %s", b.Add(*f.CarID)))
	}
	if f.Type != nil {
		clauses = append(clauses, fmt.Sprintf("type = %s", b.Add(string(*f.Type))))
	}
	if f.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = %s", b.Add(string(*f.Status))))
	}
	if f.ExpiringWithinDays != nil {
		clauses = append(clauses, fmt.Sprintf(
			"expires_at <= NOW() + make_interval(days => %s)", b.Add(*f.ExpiringWithinDays),
		))
	}

	return clauses
}

func BuildCarInsuranceSetClauses(u model.CarInsuranceUpdate, b *ArgsBuilder) []string {
	var clauses []string

	if u.Provider != nil {
		clauses = append(clauses, fmt.Sprintf("provider = %s", b.Add(*u.Provider)))
	}
	if u.PolicyNum != nil {
		clauses = append(clauses, fmt.Sprintf("policy_num = %s", b.Add(*u.PolicyNum)))
	}
	if u.StartsAt != nil {
		clauses = append(clauses, fmt.Sprintf("starts_at = %s", b.Add(*u.StartsAt)))
	}
	if u.ExpiresAt != nil {
		clauses = append(clauses, fmt.Sprintf("expires_at = %s", b.Add(*u.ExpiresAt)))
	}
	if u.CostTenge != nil {
		clauses = append(clauses, fmt.Sprintf("cost_tenge = %s", b.Add(*u.CostTenge)))
	}
	if u.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = %s", b.Add(string(*u.Status))))
	}
	if u.Notes != nil {
		clauses = append(clauses, fmt.Sprintf("notes = %s", b.Add(*u.Notes)))
	}
	if u.ImageKeys != nil {
		clauses = append(clauses, fmt.Sprintf("image_keys = %s", b.Add(pq.StringArray(u.ImageKeys))))
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = %s", b.Add(u.UpdatedAt)))

	return clauses
}
