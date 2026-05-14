package model

import "time"

type TripStatusReading struct {
	ID         string
	TripID     string
	FromStatus TripStatus
	ToStatus   TripStatus
	ActorType  ActorType
	ActorID    *string
	Reason     *string
	ChangedAt  time.Time
}

// TripStatusReadingCreate is the repo-layer input for inserting a status transition record.
type TripStatusReadingCreate struct {
	TripID     string
	FromStatus TripStatus
	ToStatus   TripStatus
	ActorType  ActorType
	ActorID    *string
	Reason     *string
	ChangedAt  time.Time
}

// TripStatusReadingFilter is used by GetTripStatusHistory.
type TripStatusReadingFilter struct {
	TripID     string
	From       *time.Time
	To         *time.Time
	Pagination *Pagination
}
