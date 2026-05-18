package dto

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"carsharing/trip-service/internal/model"
)

const TripColumns = `
	id, booking_id, user_id, car_id, status,
	started_at, start_latitude, start_longitude, start_mileage_km, start_fuel_level,
	ended_at, end_latitude, end_longitude, end_mileage_km, end_fuel_level,
	distance_traveled_km, duration_seconds, final_cost_tenge, cancel_reason,
	created_at, updated_at`

type scanner interface {
	Scan(dest ...any) error
}

type tripRow struct {
	ID             string
	BookingID      string
	UserID         string
	CarID          string
	Status         string
	StartedAt      time.Time
	StartLatitude  float64
	StartLongitude float64
	StartMileageKM int64
	StartFuelLevel sql.NullFloat64
	EndedAt        sql.NullTime
	EndLatitude    sql.NullFloat64
	EndLongitude   sql.NullFloat64
	EndMileageKM   sql.NullInt64
	EndFuelLevel   sql.NullFloat64
	DistanceKM     sql.NullFloat64
	DurationSecs   sql.NullInt64
	FinalCost      sql.NullInt32
	CancelReason   sql.NullString
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (r tripRow) toDomain() model.Trip {
	t := model.Trip{
		ID:        r.ID,
		BookingID: r.BookingID,
		UserID:    r.UserID,
		CarID:     r.CarID,
		Status:    model.TripStatus(r.Status),
		StartedAt: r.StartedAt,
		StartLocation: model.Location{
			Latitude:  r.StartLatitude,
			Longitude: r.StartLongitude,
		},
		StartMileageKM: r.StartMileageKM,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
	if r.StartFuelLevel.Valid {
		f := float32(r.StartFuelLevel.Float64)
		t.StartFuelLevel = &f
	}
	if r.EndedAt.Valid {
		t.EndedAt = &r.EndedAt.Time
	}
	if r.EndLatitude.Valid && r.EndLongitude.Valid {
		loc := model.Location{Latitude: r.EndLatitude.Float64, Longitude: r.EndLongitude.Float64}
		t.EndLocation = &loc
	}
	if r.EndMileageKM.Valid {
		m := r.EndMileageKM.Int64
		t.EndMileageKM = &m
	}
	if r.EndFuelLevel.Valid {
		f := float32(r.EndFuelLevel.Float64)
		t.EndFuelLevel = &f
	}
	if r.DistanceKM.Valid {
		t.DistanceTraveledKM = &r.DistanceKM.Float64
	}
	if r.DurationSecs.Valid {
		t.DurationSeconds = &r.DurationSecs.Int64
	}
	if r.FinalCost.Valid {
		c := r.FinalCost.Int32
		t.FinalCostTenge = &c
	}
	if r.CancelReason.Valid {
		t.CancelReason = &r.CancelReason.String
	}
	return t
}

func ScanTrip(s scanner) (model.Trip, error) {
	var r tripRow
	err := s.Scan(
		&r.ID, &r.BookingID, &r.UserID, &r.CarID, &r.Status,
		&r.StartedAt, &r.StartLatitude, &r.StartLongitude, &r.StartMileageKM, &r.StartFuelLevel,
		&r.EndedAt, &r.EndLatitude, &r.EndLongitude, &r.EndMileageKM, &r.EndFuelLevel,
		&r.DistanceKM, &r.DurationSecs, &r.FinalCost, &r.CancelReason,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.Trip{}, err
	}
	return r.toDomain(), nil
}

func BuildTripSetClauses(u model.TripUpdate, b *ArgsBuilder) []string {
	var clauses []string
	if u.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = %s", b.Add(u.Status.String())))
	}
	if u.EndedAt != nil {
		clauses = append(clauses, fmt.Sprintf("ended_at = %s", b.Add(*u.EndedAt)))
	}
	if u.EndLocation != nil {
		clauses = append(clauses, fmt.Sprintf("end_latitude = %s", b.Add(u.EndLocation.Latitude)))
		clauses = append(clauses, fmt.Sprintf("end_longitude = %s", b.Add(u.EndLocation.Longitude)))
	}
	if u.EndMileageKM != nil {
		clauses = append(clauses, fmt.Sprintf("end_mileage_km = %s", b.Add(*u.EndMileageKM)))
	}
	if u.EndFuelLevel != nil {
		clauses = append(clauses, fmt.Sprintf("end_fuel_level = %s", b.Add(float64(*u.EndFuelLevel))))
	}
	if u.DistanceTraveledKM != nil {
		clauses = append(clauses, fmt.Sprintf("distance_traveled_km = %s", b.Add(*u.DistanceTraveledKM)))
	}
	if u.DurationSeconds != nil {
		clauses = append(clauses, fmt.Sprintf("duration_seconds = %s", b.Add(*u.DurationSeconds)))
	}
	if u.FinalCostTenge != nil {
		clauses = append(clauses, fmt.Sprintf("final_cost_tenge = %s", b.Add(*u.FinalCostTenge)))
	}
	if u.CancelReason != nil {
		clauses = append(clauses, fmt.Sprintf("cancel_reason = %s", b.Add(*u.CancelReason)))
	}
	clauses = append(clauses, fmt.Sprintf("updated_at = %s", b.Add(u.UpdatedAt)))
	return clauses
}

func BuildTripWhereClauses(f model.TripFilter, b *ArgsBuilder) []string {
	var clauses []string
	if f.UserID != nil {
		clauses = append(clauses, fmt.Sprintf("user_id = %s", b.Add(*f.UserID)))
	}
	if f.CarID != nil {
		clauses = append(clauses, fmt.Sprintf("car_id = %s", b.Add(*f.CarID)))
	}
	if f.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = %s", b.Add(f.Status.String())))
	}
	if f.StartedAfter != nil {
		clauses = append(clauses, fmt.Sprintf("started_at >= %s", b.Add(*f.StartedAfter)))
	}
	if f.StartedBefore != nil {
		clauses = append(clauses, fmt.Sprintf("started_at <= %s", b.Add(*f.StartedBefore)))
	}
	return clauses
}

func BuildTripWhereClause(clauses []string) string {
	if len(clauses) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(clauses, " AND ")
}
