package dto

import (
	"carsharing/car-service/internal/model"
	sharedmodel "carsharing/shared/model"
	"fmt"
	"time"
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
	Notes            []string
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
		Location: model.Location{
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
		clauses = append(clauses, fmt.Sprintf("notes = %s", b.Add(u.Notes)))
	}
	if u.ImageKeys != nil {
		clauses = append(clauses, fmt.Sprintf("image_keys = %s", b.Add(u.ImageKeys)))
	}
	if u.LastSeenAt != nil {
		clauses = append(clauses, fmt.Sprintf("last_seen_at = %s", b.Add(*u.LastSeenAt)))
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = %s", b.Add(u.UpdatedAt)))

	return clauses
}
