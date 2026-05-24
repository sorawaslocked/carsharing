package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type TripStatusReading struct {
	ID         string
	TripID     string
	FromStatus TripStatus
	ToStatus   TripStatus
	ActorType  sharedmodel.ActorType
	ActorID    *string
	Reason     *string
	ChangedAt  time.Time
}

type TripStatusReadingCreate struct {
	TripID     string
	FromStatus TripStatus
	ToStatus   TripStatus
	ActorType  sharedmodel.ActorType
	ActorID    *string
	Reason     *string
	ChangedAt  time.Time
}

type TripStatusReadingFilter struct {
	TripID     string
	TimeRange  *sharedmodel.TimeRange
	Pagination *sharedmodel.Pagination
}
