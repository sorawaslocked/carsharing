package model

import sharedmodel "carsharing/shared/model"

import "time"

type Trip struct {
	ID        string
	BookingID string
	UserID    string
	CarID     string
	Status    string

	StartedAt      time.Time
	StartLocation  sharedmodel.Location
	StartMileageKM int64
	StartFuelLevel *float32

	EndedAt      *time.Time
	EndLocation  *sharedmodel.Location
	EndMileageKM *int64
	EndFuelLevel *float32

	DistanceTraveledKM *float64
	DurationSeconds    *int64
	FinalCostTenge     *int32

	CancelReason *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type TripFilter struct {
	UserID        *string
	CarID         *string
	Status        *string
	StartedAfter  *time.Time
	StartedBefore *time.Time
	Pagination    *sharedmodel.Pagination
}

type TripSummary struct {
	TripID    string
	BookingID string
	StartedAt time.Time
	EndedAt   time.Time

	DurationSeconds    int64
	DistanceTraveledKM float64

	PricingSnapshot        PricingSnapshot
	BaseCostTenge          int32
	DistanceCostTenge      int32
	OvertimeCostTenge      int32
	ZoneFeeAdjustmentTenge int32
	TotalCostTenge         int32
}

type TripStatusReading struct {
	ID         string
	TripID     string
	FromStatus string
	ToStatus   string
	ActorType  string
	ActorID    *string
	Reason     *string
	ChangedAt  time.Time
}

type TripStatusReadingFilter struct {
	From       *time.Time
	To         *time.Time
	Pagination *sharedmodel.Pagination
}

type TripLiveFeed struct {
	ElapsedSeconds     int64
	CurrentCostTenge   int32
	DistanceTraveledKM float64
}
