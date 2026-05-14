package model

import "time"

type TripStatus string

const (
	TripStatusActive    TripStatus = "active"
	TripStatusCompleted TripStatus = "completed"
	TripStatusCancelled TripStatus = "cancelled"
)

func (s TripStatus) String() string {
	return string(s)
}

var validTransitions = map[TripStatus]map[TripStatus]struct{}{
	TripStatusActive: {
		TripStatusCompleted: {},
		TripStatusCancelled: {},
	},
}

func (s TripStatus) CanTransitionTo(next TripStatus) bool {
	allowed, ok := validTransitions[s]
	if !ok {
		return false
	}
	_, ok = allowed[next]
	return ok
}

type Trip struct {
	ID        string
	BookingID string
	UserID    string
	CarID     string
	Status    TripStatus

	StartedAt      time.Time
	StartLocation  Location
	StartMileageKM int64
	StartFuelLevel *float32

	EndedAt      *time.Time
	EndLocation  *Location
	EndMileageKM *int64
	EndFuelLevel *float32

	DistanceTraveledKM *float64
	DurationSeconds    *int64
	FinalCostTenge     *int32
	CancelReason       *string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type TripCreate struct {
	ID             string
	BookingID      string
	UserID         string
	CarID          string
	Status         TripStatus
	StartedAt      time.Time
	StartLocation  Location
	StartMileageKM int64
	StartFuelLevel *float32
}

type TripUpdate struct {
	Status             *TripStatus
	EndedAt            *time.Time
	EndLocation        *Location
	EndMileageKM       *int64
	EndFuelLevel       *float32
	DistanceTraveledKM *float64
	DurationSeconds    *int64
	FinalCostTenge     *int32
	CancelReason       *string
	UpdatedAt          time.Time
}

// TripFilter is used by ListTrips to narrow results.
type TripFilter struct {
	UserID        *string
	CarID         *string
	Status        *TripStatus
	StartedAfter  *time.Time
	StartedBefore *time.Time
	Pagination    *Pagination
}
