package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type Car struct {
	ID               string
	ModelID          string
	VIN              string
	LicensePlate     string
	Color            string
	YearManufactured int16

	MileageKM    int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     sharedmodel.Location

	TelemetryID string
	ZoneID      *string
	FuelStatus  string
	Status      string
	IsRetired   bool

	Notes     *string
	ImageURLs []string

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CarFilter struct {
	Brand        *string
	Model        *string
	FuelType     *string
	Transmission *string
	BodyType     *string
	Class        *string
	MinSeats     *int8

	Location     *sharedmodel.Location
	RadiusM      *int32
	MinFuelLevel *float32

	ZoneID    *string
	Status    *string
	IsRetired *bool

	Pagination *sharedmodel.Pagination
}

type CarCreate struct {
	ModelID          string
	VIN              string
	LicensePlate     string
	Color            string
	YearManufactured int16

	TelemetryID string
	ZoneID      *string

	MileageKM    *int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     *sharedmodel.Location

	Notes *string
}

type CarUpdate struct {
	ModelID      *string
	LicensePlate *string
	Color        *string

	TelemetryID *string
	ZoneID      *string

	IsRetired *bool
	Notes     *string
	ImageKeys []string
}

type CarStatusReading struct {
	ID    string
	CarID string

	FromStatus string
	ToStatus   string

	ActorType string
	ActorID   *string
	Reason    *string
	Metadata  map[string]any

	RecordedAt time.Time
}

type CarStatusReadingFilter struct {
	TimeRange  *sharedmodel.TimeRange
	Pagination *sharedmodel.Pagination
}

type CarTelemetryReading struct {
	ID           string
	CarID        string
	FuelPct      *float32
	FuelRawPct   *float32
	BatteryLevel *float32
	MileageKM    *int64
	Location     *sharedmodel.Location
	ActorType    string
	ActorID      *string
	Reason       *string
	Metadata     map[string]any
	RecordedAt   time.Time
}

type CarTelemetryReadingFilter struct {
	TimeRange  *sharedmodel.TimeRange
	Pagination *sharedmodel.Pagination
}

type CarTelemetryUpdate struct {
	MileageKM    *int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     *sharedmodel.Location

	Reason   string
	Metadata map[string]any
}

type CarStatusUpdate struct {
	Status string

	Reason   string
	Metadata map[string]any
}

type SlimCar struct {
	ID           string
	ModelID      string
	LicensePlate string
	Color        string
	Location     sharedmodel.Location
	FuelLevel    float32
	Status       string
}

type CarTelemetryEvent struct {
	FuelLevel    float32
	BatteryLevel float32
	MileageKM    int64
	Location     sharedmodel.Location
	RecordedAt   time.Time
}

type CarStatusUpdatedEvent struct {
	CarID      string
	FromStatus string
	ToStatus   string
}
