package model

import (
	"time"

	sharedmodel "carsharing/shared/model"
)

type CarCreatedEvent struct {
	CarID        string
	MileageKM    int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     sharedmodel.Location
}

type BookingCreatedEvent struct {
	BookingID string
	CarID     string
	UserID    string
	StartsAt  time.Time
	EndsAt    time.Time
}

type BookingCompletedEvent struct {
	BookingID   string
	CarID       string
	UserID      string
	CompletedAt time.Time
}

type BookingCancelledEvent struct {
	BookingID string
	CarID     string
	UserID    string
	Reason    string
}

type BookingExpiredEvent struct {
	BookingID string
	CarID     string
	UserID    string
	ExpiredAt time.Time
}

type TripStartedEvent struct {
	TripID    string
	BookingID string
	CarID     string
	UserID    string
	StartedAt time.Time
}

type TripEndedEvent struct {
	TripID    string
	BookingID string
	CarID     string
	UserID    string
	EndedAt   time.Time
}

type TripCancelledEvent struct {
	TripID      string
	BookingID   string
	CarID       string
	UserID      string
	Reason      string
	CancelledAt time.Time
}
