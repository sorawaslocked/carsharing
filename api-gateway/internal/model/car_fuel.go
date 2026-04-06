package model

import "time"

type CarFuelReading struct {
	CarID      string
	FuelPct    int
	RawPct     int
	RecordedAt time.Time
}

type CarFuelReadingFilter struct {
	CarID      *string
	From       *time.Time
	To         *time.Time
	Pagination *Pagination
}
