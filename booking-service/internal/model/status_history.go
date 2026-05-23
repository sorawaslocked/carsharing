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
	ActorType  string
	ActorID    *string
	Reason     *string
	ChangedAt  time.Time
}

type BookingStatusHistoryFilter struct {
	BookingID  string
	From       *time.Time
	To         *time.Time
	Pagination sharedmodel.Pagination
}
