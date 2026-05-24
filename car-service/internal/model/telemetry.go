package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type TelemetryReading struct {
	ID           string
	CarID        string
	Location     *sharedmodel.Location
	FuelPct      *float32
	FuelRawPct   *float32
	BatteryLevel *float32
	MileageKM    *int64
	ActorID      *string
	ActorType    sharedmodel.ActorType
	Reason       *string
	Metadata     map[string]any
	RecordedAt   time.Time
}

type TelemetryUpdate struct {
	CarID        string
	Latitude     float64
	Longitude    float64
	FuelLevel    *float32
	BatteryLevel *float32
	MileageKM    int64
	RecordedAt   time.Time
}

type TelemetryReadingFilter struct {
	CarID     string
	TimeRange *sharedmodel.TimeRange

	Pagination *sharedmodel.Pagination
}
