package dto

import (
	"fmt"
	"time"

	"carsharing/car-service/internal/model"
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
	Notes     *string
	ImageKeys []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r carInsuranceRow) toDomain() model.CarInsurance {
	return model.CarInsurance{
		ID:        r.ID,
		CarID:     r.CarID,
		Type:      model.InsuranceType(r.Type),
		Provider:  r.Provider,
		PolicyNum: r.PolicyNum,
		StartsAt:  r.StartsAt,
		ExpiresAt: r.ExpiresAt,
		CostTenge: r.CostTenge,
		Status:    model.InsuranceStatus(r.Status),
		Notes:     r.Notes,
		Images:    ImageKeysToImages(r.ImageKeys),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
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

func WhereClausesFromCarInsuranceFilter(f model.CarInsuranceFilter, args []any, n int) ([]string, []any, int) {
	var clauses []string

	if f.CarID != nil {
		n++
		args = append(args, *f.CarID)
		clauses = append(clauses, fmt.Sprintf("car_id = $%d", n))
	}
	if f.Type != nil {
		n++
		args = append(args, string(*f.Type))
		clauses = append(clauses, fmt.Sprintf("type = $%d", n))
	}
	if f.Status != nil {
		n++
		args = append(args, string(*f.Status))
		clauses = append(clauses, fmt.Sprintf("status = $%d", n))
	}
	if f.ExpiringWithinDays != nil {
		n++
		args = append(args, *f.ExpiringWithinDays)
		clauses = append(clauses, fmt.Sprintf("expires_at <= NOW() + make_interval(days => $%d)", n))
	}

	return clauses, args, n
}

func SetClausesFromCarInsuranceUpdate(update model.CarInsuranceUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	n := 0

	if update.Provider != nil {
		n++
		args = append(args, *update.Provider)
		clauses = append(clauses, fmt.Sprintf("provider = $%d", n))
	}
	if update.PolicyNum != nil {
		n++
		args = append(args, *update.PolicyNum)
		clauses = append(clauses, fmt.Sprintf("policy_num = $%d", n))
	}
	if update.StartsAt != nil {
		n++
		args = append(args, *update.StartsAt)
		clauses = append(clauses, fmt.Sprintf("starts_at = $%d", n))
	}
	if update.ExpiresAt != nil {
		n++
		args = append(args, *update.ExpiresAt)
		clauses = append(clauses, fmt.Sprintf("expires_at = $%d", n))
	}
	if update.CostTenge != nil {
		n++
		args = append(args, *update.CostTenge)
		clauses = append(clauses, fmt.Sprintf("cost_tenge = $%d", n))
	}
	if update.Status != nil {
		n++
		args = append(args, string(*update.Status))
		clauses = append(clauses, fmt.Sprintf("status = $%d", n))
	}
	if update.Notes != nil {
		n++
		args = append(args, *update.Notes)
		clauses = append(clauses, fmt.Sprintf("notes = $%d", n))
	}
	if update.ImageKeys != nil {
		n++
		args = append(args, update.ImageKeys)
		clauses = append(clauses, fmt.Sprintf("image_keys = $%d", n))
	}

	n++
	args = append(args, update.UpdatedAt)
	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", n))

	return clauses, args, n
}
