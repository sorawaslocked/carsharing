package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarModel struct {
	ID           string
	Brand        string
	Model        string
	Year         int16
	FuelType     CarFuelType
	Transmission CarTransmission
	BodyType     CarBodyType
	Class        CarClass
	Seats        int8
	EngineVolume *float32
	RangeKM      int32
	Features     []string
	Images       []sharedmodel.Image

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CarModelFilter struct {
	ID           *string
	Brand        *string
	Model        *string
	FuelType     *CarFuelType
	Transmission *CarTransmission
	BodyType     *CarBodyType
	Class        *CarClass
	MinSeats     *int8

	Pagination *sharedmodel.Pagination
}

type CarModelUpdate struct {
	Brand        *string
	Model        *string
	Year         *int16
	FuelType     *CarFuelType
	Transmission *CarTransmission
	BodyType     *CarBodyType
	Class        *CarClass
	Seats        *int8
	EngineVolume *float32
	RangeKM      *int32
	Features     []string
	ImageKeys    []string
	UpdatedAt    time.Time
}
