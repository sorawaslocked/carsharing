package model

import "time"

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
	Location     Location

	TelematicsID string
	ZoneID       *string
	FuelStatus   string
	Status       string
	IsRetired    bool

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

	Location     *Location
	RadiusM      *int32
	MinFuelLevel *float32

	ZoneID    *string
	Status    *string
	IsRetired *bool

	Pagination *Pagination
}

type CarCreate struct {
	ModelID          string
	VIN              string
	LicensePlate     string
	Color            string
	YearManufactured int16

	TelematicsID string

	Notes *string
}

type CarUpdate struct {
	ModelID      *string
	LicensePlate *string
	Color        *string

	TelematicsID *string
	ZoneID       *string

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

	ChangedAt time.Time
}

type CarStatusReadingFilter struct {
	From       *time.Time
	To         *time.Time
	Pagination *Pagination
}

type CarFuelReading struct {
	ID         string
	CarID      string
	FuelPct    float32
	RawPct     float32
	ActorType  string
	ActorID    *string
	Reason     *string
	Metadata   map[string]any
	RecordedAt time.Time
}

type CarFuelReadingFilter struct {
	From       *time.Time
	To         *time.Time
	Pagination *Pagination
}

type CarLocationReading struct {
	ID         string
	CarID      string
	Location   Location
	ActorType  string
	ActorID    *string
	Reason     *string
	Metadata   map[string]any
	RecordedAt time.Time
}

type CarLocationReadingFilter struct {
	From       *time.Time
	To         *time.Time
	Pagination *Pagination
}

type CarBatteryReading struct {
	ID           string
	CarID        string
	BatteryLevel float32
	ActorType    string
	ActorID      *string
	Reason       *string
	Metadata     map[string]any
	RecordedAt   time.Time
}

type CarBatteryReadingFilter struct {
	From       *time.Time
	To         *time.Time
	Pagination *Pagination
}

type CarMileageReading struct {
	ID         string
	CarID      string
	MileageKM  int64
	ActorType  string
	ActorID    *string
	Reason     *string
	Metadata   map[string]any
	RecordedAt time.Time
}

type CarMileageReadingFilter struct {
	From       *time.Time
	To         *time.Time
	Pagination *Pagination
}

type CarTelemetryUpdate struct {
	MileageKM    *int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     *Location

	Reason   string
	Metadata map[string]any
}

type CarStatusUpdate struct {
	Status string

	Reason   string
	Metadata map[string]any
}
