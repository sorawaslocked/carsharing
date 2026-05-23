package validation

import sharedmodel "carsharing/shared/model"

type CarFilter struct {
	ID             *string         `validate:"omitempty,uuid"`
	ModelFilter    *CarModelFilter `validate:"omitempty"`
	Status         *string         `validate:"omitempty,carstatus"`
	LocationFilter *LocationFilter `validate:"omitempty"`
	Pagination     *sharedmodel.Pagination
}

type CarCreate struct {
	ModelID          string   `validate:"required"`
	VIN              string   `validate:"required,min=17,max=17,alphanum"`
	LicensePlate     string   `validate:"required,min=1,max=20"`
	Color            string   `validate:"required,min=1,max=50"`
	YearManufactured int16    `validate:"required,min=1886"`
	TelemetryID      string   `validate:"required"`
	MileageKM        *int64   `validate:"omitempty,min=0"`
	FuelLevel        *float32 `validate:"omitempty,min=0,max=100"`
	BatteryLevel     *float32 `validate:"omitempty,min=0,max=100"`
	Notes            *string  `validate:"omitempty,min=1,max=500"`
}

type CarUpdate struct {
	ModelID      *string  `validate:"omitempty"`
	LicensePlate *string  `validate:"omitempty,min=1,max=20"`
	Color        *string  `validate:"omitempty,min=1,max=50"`
	TelemetryID  *string  `validate:"omitempty"`
	ZoneID       *string  `validate:"omitempty,uuid"`
	IsRetired    *bool    `validate:"omitempty"`
	Notes        *string  `validate:"omitempty,min=1,max=500"`
	ImageKeys    []string `validate:"omitempty,max=20,dive,min=1"`
}

type CarStatusUpdate struct {
	Status string `validate:"required,carstatus"`
}
