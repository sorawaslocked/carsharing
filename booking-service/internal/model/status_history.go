package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type BookingStatusReading struct {
	ID         string
	BookingID  string
	FromStatus string
	ToStatus   string
	ActorType  sharedmodel.ActorType
	ActorID    *string
	Reason     *string
	ChangedAt  time.Time
}

type BookingStatusHistoryFilter struct {
	BookingID  string
	TimeRange  *sharedmodel.TimeRange
	Pagination sharedmodel.Pagination
}
