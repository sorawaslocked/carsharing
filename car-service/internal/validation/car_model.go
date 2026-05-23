package validation

import sharedvalidation "carsharing/shared/validation"

type CarModelFilter struct {
	ID           *string `validate:"omitempty,uuid"`
	Brand        *string `validate:"omitempty,min=1,max=100"`
	Model        *string `validate:"omitempty,min=1,max=100"`
	FuelType     *string `validate:"omitempty,carfueltype"`
	Transmission *string `validate:"omitempty,cartransmission"`
	BodyType     *string `validate:"omitempty,carbodytype"`
	Class        *string `validate:"omitempty,carclass"`
	MinSeats     *int8   `validate:"omitempty,min=1,max=9"`
	Pagination   *sharedvalidation.Pagination
}

type CarModelCreate struct {
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

type CarModelUpdate struct {
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
	ImageKeys    []string `validate:"omitempty,max=20,dive,min=1"`
}
