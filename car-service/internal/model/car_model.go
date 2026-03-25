package model

import "time"

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

	Pagination
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
	UpdatedAt    time.Time
}

type CarModelFilterInput struct {
	ID           *string `validate:"required_without_all=Brand Model FuelType Transmission BodyType Class MinSeats"`
	Brand        *string `validate:"omitempty,min=1,max=100"`
	Model        *string `validate:"omitempty,min=1,max=100"`
	FuelType     *string `validate:"omitempty,carfueltype"`
	Transmission *string `validate:"omitempty,cartransmission"`
	BodyType     *string `validate:"omitempty,carbodytype"`
	Class        *string `validate:"omitempty,carclass"`
	MinSeats     *int8   `validate:"omitempty,min=1,max=9"`

	PaginationInput
}

type CarModelCreateInput struct {
	Brand        string   `validate:"required,min=1,max=100"`
	Model        string   `validate:"required,min=1,max=100"`
	Year         int16    `validate:"required,min=1886"`
	FuelType     string   `validate:"required,carfueltype"`
	Transmission string   `validate:"required,cartransmission"`
	BodyType     string   `validate:"required,carbodytype"`
	Class        string   `validate:"required,carclass"`
	Seats        int8     `validate:"required,min=1,max=9"`
	EngineVolume *float32 `validate:"omitempty,min=0.1,max=28.5"`
	RangeKM      int32    `validate:"min=0"`
	Features     []string `validate:"omitempty,max=50,dive,min=1,max=50"`
}

type CarModelUpdateInput struct {
	Brand        *string  `validate:"omitempty,min=1,max=100"`
	Model        *string  `validate:"omitempty,min=1,max=100"`
	Year         *int16   `validate:"omitempty,min=1886"`
	FuelType     *string  `validate:"omitempty,carfueltype"`
	Transmission *string  `validate:"omitempty,cartransmission"`
	BodyType     *string  `validate:"omitempty,carbodytype"`
	Class        *string  `validate:"omitempty,carclass"`
	Seats        *int8    `validate:"omitempty,min=1,max=9"`
	EngineVolume *float32 `validate:"omitempty,min=0.1,max=28.5"`
	RangeKM      *int32   `validate:"omitempty,min=0"`
	Features     []string `validate:"omitempty,max=50,dive,min=1,max=50"`
}
