package model

import "time"

type CarModel struct {
	ID               string
	Brand            string
	Model            string
	Year             int16
	FuelType         string
	Transmission     string
	BodyType         string
	Class            string
	Seats            int8
	EngineVolume     *float32
	RangeKM          int32
	Features         []string
	ImageStorageUrls []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CarModelFilter struct {
	Brand        *string
	Model        *string
	FuelType     *string
	Transmission *string
	BodyType     *string
	Class        *string
	MinSeats     *int8

	Pagination *Pagination
}

type CarModelCreate struct {
	Brand            string
	Model            string
	Year             int16
	FuelType         string
	Transmission     string
	BodyType         string
	Class            string
	Seats            int8
	EngineVolume     *float32
	RangeKM          int32
	Features         []string
	ImageStorageKeys []string
}

type CarModelUpdate struct {
	Brand            *string
	Model            *string
	Year             *int16
	FuelType         *string
	Transmission     *string
	BodyType         *string
	Class            *string
	Seats            *int8
	EngineVolume     *float32
	RangeKM          *int32
	Features         []string
	ImageStorageKeys []string
}
