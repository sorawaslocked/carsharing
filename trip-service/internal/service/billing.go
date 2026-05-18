package service

import "carsharing/trip-service/internal/model"

// calculateCosts breaks down the total fare into base, distance, and overtime components.
// durationSeconds is the elapsed trip time; distanceKM is the odometer delta.
func calculateCosts(ps model.PricingSnapshot, committedPeriods *int32, durationSeconds int64, distanceKM float64) (base, distance, overtime int32) {
	durationMinutes := int32(durationSeconds / 60)

	freeMinutes := int32(0)
	if ps.FreeMinutes != nil {
		freeMinutes = *ps.FreeMinutes
	}

	billableMinutes := durationMinutes - freeMinutes
	if billableMinutes < 0 {
		billableMinutes = 0
	}

	base = ps.RateTenge * billableMinutes
	if ps.MinChargeTenge != nil && base < *ps.MinChargeTenge {
		base = *ps.MinChargeTenge
	}

	if ps.RatePerKmTenge != nil {
		distance = int32(distanceKM * float64(*ps.RatePerKmTenge))
	}

	if committedPeriods != nil && ps.OvertimePolicy != nil && ps.OvertimeRateTenge != nil {
		overtimeMinutes := durationMinutes - *committedPeriods
		if overtimeMinutes > 0 {
			overtime = overtimeMinutes * *ps.OvertimeRateTenge
		}
	}

	return
}

func calculateCurrentCost(ps model.PricingSnapshot, committedPeriods *int32, elapsedSeconds int64, distanceKM float64) int32 {
	base, dist, overtime := calculateCosts(ps, committedPeriods, elapsedSeconds, distanceKM)
	return base + dist + overtime
}
