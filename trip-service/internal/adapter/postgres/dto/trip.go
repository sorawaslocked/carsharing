package dto

import (
	"fmt"
	"strings"
	"time"

	sharedmodel "carsharing/shared/model"
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
	StartFuelLevel *float32
	EndedAt        *time.Time
	EndLatitude    *float64
	EndLongitude   *float64
	EndMileageKM   *int64
	EndFuelLevel   *float32
	DistanceKM     *float64
	DurationSecs   *int64
	FinalCost      *int32
	CancelReason   *string
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
		StartLocation: sharedmodel.Location{
			Latitude:  r.StartLatitude,
			Longitude: r.StartLongitude,
		},
		StartMileageKM: r.StartMileageKM,
		StartFuelLevel: r.StartFuelLevel,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
	t.EndedAt = r.EndedAt
	if r.EndLatitude != nil && r.EndLongitude != nil {
		loc := sharedmodel.Location{Latitude: *r.EndLatitude, Longitude: *r.EndLongitude}
		t.EndLocation = &loc
	}
	t.EndMileageKM = r.EndMileageKM
	t.EndFuelLevel = r.EndFuelLevel
	t.DistanceTraveledKM = r.DistanceKM
	t.DurationSeconds = r.DurationSecs
	t.FinalCostTenge = r.FinalCost
	t.CancelReason = r.CancelReason
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
		clauses = append(clauses, fmt.Sprintf("end_fuel_level = %s", b.Add(*u.EndFuelLevel)))
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
