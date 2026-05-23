package dto

import (
	"fmt"
	"time"

	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"
)

type carRow struct {
	ID               string
	ModelID          string
	VIN              string
	LicensePlate     string
	Color            string
	YearManufactured int16
	Status           string
	MileageKM        int64
	FuelLevel        *float32
	BatteryLevel     *float32
	Latitude         float64
	Longitude        float64
	Notes            *string
	ImageKeys        []string
	LastSeenAt       time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (r carRow) toDomain() model.Car {
	return model.Car{
		ID:               r.ID,
		ModelID:          r.ModelID,
		VIN:              r.VIN,
		LicensePlate:     r.LicensePlate,
		Color:            r.Color,
		YearManufactured: r.YearManufactured,
		Status:           model.CarStatus(r.Status),
		MileageKM:        r.MileageKM,
		FuelLevel:        r.FuelLevel,
		BatteryLevel:     r.BatteryLevel,
		Location: sharedmodel.Location{
			Latitude:  r.Latitude,
			Longitude: r.Longitude,
		},
		Notes:      r.Notes,
		Images:     ImageKeysToImages(r.ImageKeys),
		LastSeenAt: r.LastSeenAt,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func ScanCarRow(s scanner) (model.Car, error) {
	var r carRow

	err := s.Scan(
		&r.ID, &r.ModelID, &r.VIN, &r.LicensePlate, &r.Color,
		&r.YearManufactured, &r.Status, &r.MileageKM,
		&r.FuelLevel, &r.BatteryLevel,
		&r.Latitude, &r.Longitude,
		&r.Notes, &r.ImageKeys, &r.LastSeenAt, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return model.Car{}, err
	}

	return r.toDomain(), nil
}

func SetClausesFromCarUpdate(update model.CarUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	n := 0

	if update.ModelID != nil {
		n++
		args = append(args, *update.ModelID)
		clauses = append(clauses, fmt.Sprintf("model_id = $%d", n))
	}
	if update.LicensePlate != nil {
		n++
		args = append(args, *update.LicensePlate)
		clauses = append(clauses, fmt.Sprintf("license_plate = $%d", n))
	}
	if update.Color != nil {
		n++
		args = append(args, *update.Color)
		clauses = append(clauses, fmt.Sprintf("color = $%d", n))
	}
	if update.MileageKM != nil {
		n++
		args = append(args, *update.MileageKM)
		clauses = append(clauses, fmt.Sprintf("mileage_km = $%d", n))
	}
	if update.FuelLevel != nil {
		n++
		args = append(args, *update.FuelLevel)
		clauses = append(clauses, fmt.Sprintf("fuel_level = $%d", n))
	}
	if update.BatteryLevel != nil {
		n++
		args = append(args, *update.BatteryLevel)
		clauses = append(clauses, fmt.Sprintf("battery_level = $%d", n))
	}
	if update.Location != nil {
		n++
		args = append(args, update.Location.Latitude)
		clauses = append(clauses, fmt.Sprintf("latitude = $%d", n))
		n++
		args = append(args, update.Location.Longitude)
		clauses = append(clauses, fmt.Sprintf("longitude = $%d", n))
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
	if update.LastSeenAt != nil {
		n++
		args = append(args, *update.LastSeenAt)
		clauses = append(clauses, fmt.Sprintf("last_seen_at = $%d", n))
	}

	n++
	args = append(args, update.UpdatedAt)
	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", n))

	return clauses, args, n
}
