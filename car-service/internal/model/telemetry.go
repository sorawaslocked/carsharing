package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarTelemetryEvent struct {
	ID           string
	CarID        string
	Location     *sharedmodel.Location
	FuelPct      *float32
	FuelRawPct   *float32
	BatteryLevel *float32
	MileageKM    *int64
	ActorID      *string
	ActorType    string
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

type CarTelemetryUpdateInput struct {
	MileageKM    int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     *sharedmodel.Location
}

type TelemetryEventFilter struct {
	CarID *string
	From  *time.Time
	To    *time.Time

	Pagination *sharedmodel.Pagination
}
