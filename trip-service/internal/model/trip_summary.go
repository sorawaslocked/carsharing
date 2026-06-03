package model

import "time"

type PricingSnapshot struct {
	RateTenge         int32
	RatePerKmTenge    *int32
	FreeMinutes       *int32
	MinChargeTenge    *int32
	OvertimePolicy    *string
	OvertimeRateTenge *int32
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

// TripSummaryCreate is the repo-layer input for persisting a billing summary when a trip ends.
type TripSummaryCreate struct {
	TripID                 string
	BookingID              string
	StartedAt              time.Time
	EndedAt                time.Time
	DurationSeconds        int64
	DistanceTraveledKM     float64
	PricingSnapshot        PricingSnapshot
	BaseCostTenge          int32
	DistanceCostTenge      int32
	OvertimeCostTenge      int32
	ZoneFeeAdjustmentTenge int32
	TotalCostTenge         int32
}
