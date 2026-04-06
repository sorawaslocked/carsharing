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
	FuelStatus   string
	ZoneID       *string

	Status           string
	Notes            *string
	ImageStorageUrls []string

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

	Latitude     *float64
	Longitude    *float64
	RadiusM      *int32
	ZoneID       *string
	MinFuelLevel *float32

	Status *string

	Pagination *Pagination
}

type CarCreate struct {
	ModelID          string
	VIN              string
	LicensePlate     string
	Color            string
	YearManufactured int16

	MileageKM    int64
	FuelLevel    *float32
	BatteryLevel *float32
	Latitude     float64
	Longitude    float64
	TelematicsID string

	Notes            *string
	ImageStorageKeys []string
}

type CarUpdate struct {
	ModelID      *string
	LicensePlate *string
	Color        *string

	MileageKM    *int64
	FuelLevel    *float32
	BatteryLevel *float32
	Latitude     *float64
	Longitude    *float64

	TelematicsID *string
	ZoneID       *string

	Status           *string
	Notes            *string
	ImageStorageKeys []string
}
