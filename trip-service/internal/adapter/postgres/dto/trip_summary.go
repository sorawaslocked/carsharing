package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"carsharing/trip-service/internal/model"
)

const TripSummaryColumns = `
	trip_id, booking_id, started_at, ended_at,
	duration_seconds, distance_traveled_km, pricing_snapshot,
	base_cost_tenge, distance_cost_tenge, overtime_cost_tenge, zone_fee_adjustment_tenge, total_cost_tenge`

type pricingSnapshotJSON struct {
	RateTenge         int32   `json:"rate_tenge"`
	RatePerKmTenge    *int32  `json:"rate_per_km_tenge,omitempty"`
	FreeMinutes       *int32  `json:"free_minutes,omitempty"`
	MinChargeTenge    *int32  `json:"min_charge_tenge,omitempty"`
	OvertimePolicy    *string `json:"overtime_policy,omitempty"`
	OvertimeRateTenge *int32  `json:"overtime_rate_tenge,omitempty"`
}

func MarshalPricingSnapshot(ps model.PricingSnapshot) ([]byte, error) {
	return json.Marshal(pricingSnapshotJSON{
		RateTenge:         ps.RateTenge,
		RatePerKmTenge:    ps.RatePerKmTenge,
		FreeMinutes:       ps.FreeMinutes,
		MinChargeTenge:    ps.MinChargeTenge,
		OvertimePolicy:    ps.OvertimePolicy,
		OvertimeRateTenge: ps.OvertimeRateTenge,
	})
}

func unmarshalPricingSnapshot(data []byte) (model.PricingSnapshot, error) {
	var j pricingSnapshotJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return model.PricingSnapshot{}, fmt.Errorf("unmarshal pricing snapshot: %w", err)
	}
	return model.PricingSnapshot{
		RateTenge:         j.RateTenge,
		RatePerKmTenge:    j.RatePerKmTenge,
		FreeMinutes:       j.FreeMinutes,
		MinChargeTenge:    j.MinChargeTenge,
		OvertimePolicy:    j.OvertimePolicy,
		OvertimeRateTenge: j.OvertimeRateTenge,
	}, nil
}

type tripSummaryRow struct {
	TripID                 string
	BookingID              string
	StartedAt              time.Time
	EndedAt                time.Time
	DurationSeconds        int64
	DistanceTraveledKM     float64
	PricingSnapshot        []byte
	BaseCostTenge          int32
	DistanceCostTenge      int32
	OvertimeCostTenge      int32
	ZoneFeeAdjustmentTenge int32
	TotalCostTenge         int32
}

func ScanTripSummary(s scanner) (model.TripSummary, error) {
	var r tripSummaryRow
	err := s.Scan(
		&r.TripID, &r.BookingID, &r.StartedAt, &r.EndedAt,
		&r.DurationSeconds, &r.DistanceTraveledKM, &r.PricingSnapshot,
		&r.BaseCostTenge, &r.DistanceCostTenge, &r.OvertimeCostTenge, &r.ZoneFeeAdjustmentTenge, &r.TotalCostTenge,
	)
	if err != nil {
		return model.TripSummary{}, err
	}

	ps, err := unmarshalPricingSnapshot(r.PricingSnapshot)
	if err != nil {
		return model.TripSummary{}, err
	}

	return model.TripSummary{
		TripID:                 r.TripID,
		BookingID:              r.BookingID,
		StartedAt:              r.StartedAt,
		EndedAt:                r.EndedAt,
		DurationSeconds:        r.DurationSeconds,
		DistanceTraveledKM:     r.DistanceTraveledKM,
		PricingSnapshot:        ps,
		BaseCostTenge:          r.BaseCostTenge,
		DistanceCostTenge:      r.DistanceCostTenge,
		OvertimeCostTenge:      r.OvertimeCostTenge,
		ZoneFeeAdjustmentTenge: r.ZoneFeeAdjustmentTenge,
		TotalCostTenge:         r.TotalCostTenge,
	}, nil
}
