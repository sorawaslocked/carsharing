package dto

import (
	"car-rental-car-service/internal/model"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type carRow struct {
	ID               string          `db:"id"`
	ModelID          string          `db:"model_id"`
	VIN              string          `db:"vin"`
	LicensePlate     string          `db:"license_plate"`
	Color            string          `db:"color"`
	YearManufactured int16           `db:"year_manufactured"`
	Status           string          `db:"status"`
	MileageKM        int64           `db:"mileage_km"`
	FuelLevel        sql.NullFloat64 `db:"fuel_level"`
	BatteryLevel     sql.NullFloat64 `db:"battery_level"`
	Latitude         float64         `db:"latitude"`
	Longitude        float64         `db:"longitude"`
	Notes            pq.StringArray  `db:"notes"`
	LastSeenAt       time.Time       `db:"last_seen_at"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`
}

func (r carRow) toDomain() model.Car {
	c := model.Car{
		ID:               r.ID,
		ModelID:          r.ModelID,
		VIN:              r.VIN,
		LicensePlate:     r.LicensePlate,
		Color:            r.Color,
		YearManufactured: r.YearManufactured,
		Status:           model.CarStatus(r.Status),
		MileageKM:        r.MileageKM,
		Location: model.Location{
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
		},
		Notes:      []string(r.Notes),
		LastSeenAt: r.LastSeenAt,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
	if r.FuelLevel.Valid {
		c.FuelLevel = new(float32(r.FuelLevel.Float64))
	}
	if r.BatteryLevel.Valid {
		c.BatteryLevel = new(float32(r.BatteryLevel.Float64))
	}

	return c
}

func ScanCarRow(s scanner) (model.Car, error) {
	var r carRow

	err := s.Scan(
		&r.ID, &r.ModelID, &r.VIN, &r.LicensePlate, &r.Color,
		&r.YearManufactured, &r.Status, &r.MileageKM,
		&r.FuelLevel, &r.BatteryLevel,
		&r.Latitude, &r.Longitude,
		&r.Notes, &r.LastSeenAt, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.Car{}, err
	}

	return r.toDomain(), nil
}

func BuildCarSetClauses(u model.CarUpdate, b *ArgsBuilder) []string {
	var clauses []string

	if u.ModelID != nil {
		clauses = append(clauses, fmt.Sprintf("model_id = %s", b.Add(*u.ModelID)))
	}
	if u.LicensePlate != nil {
		clauses = append(clauses, fmt.Sprintf("license_plate = %s", b.Add(*u.LicensePlate)))
	}
	if u.Color != nil {
		clauses = append(clauses, fmt.Sprintf("color = %s", b.Add(*u.Color)))
	}
	if u.MileageKM != nil {
		clauses = append(clauses, fmt.Sprintf("mileage_km = %s", b.Add(*u.MileageKM)))
	}
	if u.FuelLevel != nil {
		clauses = append(clauses, fmt.Sprintf("fuel_level = %s", b.Add(*u.FuelLevel)))
	}
	if u.BatteryLevel != nil {
		clauses = append(clauses, fmt.Sprintf("battery_level = %s", b.Add(*u.BatteryLevel)))
	}
	if u.Location != nil {
		clauses = append(clauses, fmt.Sprintf("latitude = %s", b.Add(u.Location.Latitude)))
		clauses = append(clauses, fmt.Sprintf("longitude = %s", b.Add(u.Location.Longitude)))
	}
	if u.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = %s", b.Add(string(*u.Status))))
	}
	if u.Notes != nil {
		clauses = append(clauses, fmt.Sprintf("notes = %s", b.Add(pq.StringArray(u.Notes))))
	}
	if u.LastSeenAt != nil {
		clauses = append(clauses, fmt.Sprintf("last_seen_at = %s", b.Add(*u.LastSeenAt)))
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = %s", b.Add(u.UpdatedAt)))

	return clauses
}
