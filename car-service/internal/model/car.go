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
	Location     Location

	Status CarStatus
	Notes  []string
	Images []sharedmodel.Image

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CarFilter struct {
	ID             *string
	ModelFilter    *CarModelFilter
	Status         *CarStatus
	LocationFilter *LocationFilter

	Pagination *sharedmodel.Pagination
}

type CarUpdate struct {
	ModelID      *string
	LicensePlate *string
	Color        *string

	MileageKM    *int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     *Location

	Status    *CarStatus
	Notes     []string
	ImageKeys []string

	LastSeenAt *time.Time
	UpdatedAt  time.Time
}
