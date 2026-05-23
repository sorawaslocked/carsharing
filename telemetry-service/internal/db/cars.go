package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type CarTelemetry struct {
	Latitude     float64
	Longitude    float64
	FuelLevel    *float32
	BatteryLevel *float32
	MileageKm    int64
}

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository(dsn string) (*CarRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return &CarRepository{db: db}, nil
}

func (r *CarRepository) Close() error {
	return r.db.Close()
}

func (r *CarRepository) GetCarTelemetry(carID string) (*CarTelemetry, error) {
	const q = `
		SELECT latitude, longitude, fuel_level, battery_level, mileage_km
		FROM cars
		WHERE id = $1`

	var (
		t            CarTelemetry
		fuelLevel    sql.NullFloat64
		batteryLevel sql.NullFloat64
	)
	err := r.db.QueryRow(q, carID).Scan(
		&t.Latitude,
		&t.Longitude,
		&fuelLevel,
		&batteryLevel,
		&t.MileageKm,
	)
	if err != nil {
		return nil, err
	}

	if fuelLevel.Valid {
		v := float32(fuelLevel.Float64)
		t.FuelLevel = &v
	}
	if batteryLevel.Valid {
		v := float32(batteryLevel.Float64)
		t.BatteryLevel = &v
	}

	return &t, nil
}
