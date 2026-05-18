package dto

import (
	"fmt"
	"time"

	"github.com/sorawaslocked/car-rental-car-service/internal/model"
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

func BuildZoneWhereClauses(f model.ZoneFilter, b *ArgsBuilder) []string {
	var clauses []string

	if f.Type != nil {
		clauses = append(clauses, fmt.Sprintf("type = %s", b.Add(string(*f.Type))))
	}
	if f.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("is_active = %s", b.Add(*f.IsActive)))
	}

	return clauses
}

func BuildZoneSetClauses(u model.ZoneUpdate, b *ArgsBuilder) []string {
	var clauses []string

	if u.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = %s", b.Add(*u.Name)))
	}
	if u.Type != nil {
		clauses = append(clauses, fmt.Sprintf("type = %s", b.Add(string(*u.Type))))
	}
	if u.BoundaryGeoJSON != nil {
		clauses = append(clauses, fmt.Sprintf("boundary_geo_json = %s", b.Add(*u.BoundaryGeoJSON)))
	}
	if u.FeeAdjustment != nil {
		clauses = append(clauses, fmt.Sprintf("fee_adjustment = %s", b.Add(*u.FeeAdjustment)))
	}
	if u.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("is_active = %s", b.Add(*u.IsActive)))
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = %s", b.Add(u.UpdatedAt)))

	return clauses
}
