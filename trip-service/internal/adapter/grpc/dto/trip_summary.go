package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	basepb "github.com/sorawaslocked/car-rental-protos/gen/base"
	basetripmb "github.com/sorawaslocked/car-rental-protos/gen/base/trip"

	"carsharing/trip-service/internal/model"
)

func TripSummaryToProto(s model.TripSummary) *basetripmb.TripSummary {
	return &basetripmb.TripSummary{
		TripId:             s.TripID,
		BookingId:          s.BookingID,
		StartedAt:          timestamppb.New(s.StartedAt),
		EndedAt:            timestamppb.New(s.EndedAt),
		DurationSeconds:    s.DurationSeconds,
		DistanceTraveledKm: s.DistanceTraveledKM,
		PricingSnapshot:    PricingSnapshotToProto(s.PricingSnapshot),
		BaseCostTenge:      s.BaseCostTenge,
		DistanceCostTenge:  s.DistanceCostTenge,
		OvertimeCostTenge:  s.OvertimeCostTenge,
		TotalCostTenge:     s.TotalCostTenge,
	}
}

func PricingSnapshotToProto(ps model.PricingSnapshot) *basepb.PricingSnapshot {
	return &basepb.PricingSnapshot{
		RateTenge:         ps.RateTenge,
		RatePerKmTenge:    ps.RatePerKmTenge,
		FreeMinutes:       ps.FreeMinutes,
		MinChargeTenge:    ps.MinChargeTenge,
		OvertimePolicy:    ps.OvertimePolicy,
		OvertimeRateTenge: ps.OvertimeRateTenge,
	}
}

func PricingSnapshotFromProto(ps *basepb.PricingSnapshot) model.PricingSnapshot {
	if ps == nil {
		return model.PricingSnapshot{}
	}
	return model.PricingSnapshot{
		RateTenge:         ps.RateTenge,
		RatePerKmTenge:    ps.RatePerKmTenge,
		FreeMinutes:       ps.FreeMinutes,
		MinChargeTenge:    ps.MinChargeTenge,
		OvertimePolicy:    ps.OvertimePolicy,
		OvertimeRateTenge: ps.OvertimeRateTenge,
	}
}
