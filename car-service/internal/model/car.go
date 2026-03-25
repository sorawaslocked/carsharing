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

	Status CarStatus
	Notes  []string

	LastSeenAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CarFilter struct {
	ID             *string
	ModelFilter    *CarModelFilter
	Status         *CarStatus
	LocationFilter *LocationFilter

	Pagination
}

type CarUpdate struct {
	ModelID      *string
	LicensePlate *string
	Color        *string

	MileageKM    *int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     *Location

	Status *CarStatus
	Notes  []string

	LastSeenAt *time.Time
	UpdatedAt  time.Time
}

type CarFilterInput struct {
	ID             *string              `validate:"omitempty,uuid,required_without_all=ModelFilter Status LocationFilter"`
	ModelFilter    *CarModelFilterInput `validate:"omitempty"`
	Status         *string              `validate:"omitempty,carstatus"`
	LocationFilter *LocationFilter      `validate:"omitempty"`

	PaginationInput
}

type CarCreateInput struct {
	ModelID          string   `validate:"required"`
	VIN              string   `validate:"required,min=17,max=17,alphanum"`
	LicensePlate     string   `validate:"required,min=1,max=20"`
	Color            string   `validate:"required,min=1,max=50"`
	YearManufactured int16    `validate:"required,min=1886"`
	MileageKM        int64    `validate:"min=0"`
	FuelLevel        *float32 `validate:"omitempty,min=0,max=100"`
	BatteryLevel     *float32 `validate:"omitempty,min=0,max=100"`
	Notes            []string `validate:"omitempty,max=20,dive,min=1,max=500"`
}

type CarUpdateInput struct {
	ModelID      *string  `validate:"omitempty"`
	LicensePlate *string  `validate:"omitempty,min=1,max=20"`
	Color        *string  `validate:"omitempty,min=1,max=50"`
	Notes        []string `validate:"omitempty,max=20,dive,min=1,max=500"`
}

type CarStatusUpdateInput struct {
	Status string `validate:"required,carstatus"`
}
