package dto

import (
	"fmt"
	"time"

	"carsharing/car-service/internal/model"
)

type zoneRow struct {
	ID              string
	Name            string
	Type            string
	BoundaryGeoJSON string
	FeeAdjustment   int32
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (r zoneRow) toDomain() model.Zone {
	return model.Zone{
		ID:              r.ID,
		Name:            r.Name,
		Type:            model.ZoneType(r.Type),
		BoundaryGeoJSON: r.BoundaryGeoJSON,
		FeeAdjustment:   r.FeeAdjustment,
		IsActive:        r.IsActive,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func ScanZoneRow(s scanner) (model.Zone, error) {
	var r zoneRow

	err := s.Scan(
		&r.ID, &r.Name, &r.Type, &r.BoundaryGeoJSON,
		&r.FeeAdjustment, &r.IsActive, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.Zone{}, err
	}

	return r.toDomain(), nil
}

func WhereClausesFromZoneFilter(f model.ZoneFilter, args []any, n int) ([]string, []any, int) {
	var clauses []string

	if f.Type != nil {
		n++
		args = append(args, string(*f.Type))
		clauses = append(clauses, fmt.Sprintf("type = $%d", n))
	}
	if f.IsActive != nil {
		n++
		args = append(args, *f.IsActive)
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", n))
	}

	return clauses, args, n
}

func SetClausesFromZoneUpdate(update model.ZoneUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	n := 0

	if update.Name != nil {
		n++
		args = append(args, *update.Name)
		clauses = append(clauses, fmt.Sprintf("name = $%d", n))
	}
	if update.Type != nil {
		n++
		args = append(args, string(*update.Type))
		clauses = append(clauses, fmt.Sprintf("type = $%d", n))
	}
	if update.BoundaryGeoJSON != nil {
		n++
		args = append(args, *update.BoundaryGeoJSON)
		clauses = append(clauses, fmt.Sprintf("boundary_geo_json = $%d", n))
	}
	if update.FeeAdjustment != nil {
		n++
		args = append(args, *update.FeeAdjustment)
		clauses = append(clauses, fmt.Sprintf("fee_adjustment = $%d", n))
	}
	if update.IsActive != nil {
		n++
		args = append(args, *update.IsActive)
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", n))
	}

	n++
	args = append(args, update.UpdatedAt)
	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", n))

	return clauses, args, n
}
