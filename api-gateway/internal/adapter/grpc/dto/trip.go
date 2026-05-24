package dto

import (
	"carsharing/api-gateway/internal/model"
	basetripb "carsharing/protos/gen/base/trip"
)

func TripFromProto(t *basetripb.Trip) model.Trip {
	trip := model.Trip{
		ID:             t.GetId(),
		BookingID:      t.GetBookingId(),
		UserID:         t.GetUserId(),
		CarID:          t.GetCarId(),
		Status:         t.GetStatus(),
		StartMileageKM: t.GetStartMileageKm(),
		StartLocation:  LocationFromProto(t.GetStartLocation()),
		CancelReason:   t.CancelReason,
	}

	if t.StartFuelLevel != nil {
		v := t.GetStartFuelLevel()
		trip.StartFuelLevel = &v
	}
	if t.GetStartedAt() != nil {
		trip.StartedAt = t.GetStartedAt().AsTime()
	}
	if t.GetEndedAt() != nil {
		v := t.GetEndedAt().AsTime()
		trip.EndedAt = &v
	}
	if t.EndLocation != nil {
		l := LocationFromProto(t.EndLocation)
		trip.EndLocation = &l
	}
	if t.EndMileageKm != nil {
		v := t.GetEndMileageKm()
		trip.EndMileageKM = &v
	}
	if t.EndFuelLevel != nil {
		v := t.GetEndFuelLevel()
		trip.EndFuelLevel = &v
	}
	if t.DistanceTraveledKm != nil {
		v := t.GetDistanceTraveledKm()
		trip.DistanceTraveledKM = &v
	}
	if t.DurationSeconds != nil {
		v := t.GetDurationSeconds()
		trip.DurationSeconds = &v
	}
	if t.FinalCostTenge != nil {
		v := t.GetFinalCostTenge()
		trip.FinalCostTenge = &v
	}
	if t.GetCreatedAt() != nil {
		trip.CreatedAt = t.GetCreatedAt().AsTime()
	}
	if t.GetUpdatedAt() != nil {
		trip.UpdatedAt = t.GetUpdatedAt().AsTime()
	}

	return trip
}

func TripSummaryFromProto(s *basetripb.TripSummary) model.TripSummary {
	summary := model.TripSummary{
		TripID:             s.GetTripId(),
		BookingID:          s.GetBookingId(),
		DurationSeconds:    s.GetDurationSeconds(),
		DistanceTraveledKM: s.GetDistanceTraveledKm(),
		BaseCostTenge:      s.GetBaseCostTenge(),
		DistanceCostTenge:  s.GetDistanceCostTenge(),
		OvertimeCostTenge:  s.GetOvertimeCostTenge(),
		TotalCostTenge:     s.GetTotalCostTenge(),
	}

	if s.GetStartedAt() != nil {
		summary.StartedAt = s.GetStartedAt().AsTime()
	}
	if s.GetEndedAt() != nil {
		summary.EndedAt = s.GetEndedAt().AsTime()
	}
	if ps := s.GetPricingSnapshot(); ps != nil {
		summary.PricingSnapshot = PricingSnapshotFromProto(ps)
	}

	return summary
}

func TripStatusReadingFromProto(r *basetripb.TripStatusReading) model.TripStatusReading {
	reading := model.TripStatusReading{
		ID:         r.GetId(),
		TripID:     r.GetTripId(),
		FromStatus: r.GetFromStatus(),
		ToStatus:   r.GetToStatus(),
		ActorType:  r.GetActorType(),
		ActorID:    r.ActorId,
		Reason:     r.Reason,
	}

	if r.GetChangedAt() != nil {
		reading.ChangedAt = r.GetChangedAt().AsTime()
	}

	return reading
}
