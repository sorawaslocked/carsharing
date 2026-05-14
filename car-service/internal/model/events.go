package model

import "time"

type CarCreatedEvent struct {
	CarID        string
	MileageKM    int64
	FuelLevel    *float32
	BatteryLevel *float32
	Location     Location
}

type BookingCreatedEvent struct {
	BookingID string
	CarID     string
	UserID    string
	StartsAt  time.Time
	EndsAt    time.Time
}

type TripStartedEvent struct {
	TripID    string
	BookingID string
	CarID     string
	UserID    string
	StartedAt time.Time
}

type BookingCancelledEvent struct {
	BookingID string
	CarID     string
	UserID    string
	Reason    string
}

type TripEndedEvent struct {
	TripID    string
	BookingID string
	CarID     string
	UserID    string
	EndedAt   time.Time
}
