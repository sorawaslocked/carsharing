package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarStatus string

const (
	CarStatusAvailable    CarStatus = "available"
	CarStatusReserved     CarStatus = "reserved"
	CarStatusInUse        CarStatus = "in_use"
	CarStatusMaintenance  CarStatus = "maintenance"
	CarStatusOutOfService CarStatus = "out_of_service"
)

var validCarStatuses = map[CarStatus]struct{}{
	CarStatusAvailable:    {},
	CarStatusReserved:     {},
	CarStatusInUse:        {},
	CarStatusMaintenance:  {},
	CarStatusOutOfService: {},
}

func CarStatusFromString(s string) (CarStatus, bool) {
	cs := CarStatus(s)
	if _, ok := validCarStatuses[cs]; !ok {
		return "", false
	}
	return cs, true
}

func (s CarStatus) String() string {
	return string(s)
}

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
	IsRetired   bool

	Status CarStatus
	Notes  *string
	Images []sharedmodel.Image

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CarFilter struct {
	ID             *string
	ModelFilter    *CarModelFilter
	Status         *CarStatus
	IsRetired      *bool
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
	Location     *sharedmodel.Location

	TelemetryID *string
	ZoneID      *string
	IsRetired   *bool

	Status    *CarStatus
	Notes     *string
	ImageKeys []string

	LastSeenAt *time.Time
	UpdatedAt  time.Time
}
