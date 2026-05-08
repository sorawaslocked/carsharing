package model

import "time"

type Booking struct {
	ID     string
	UserID string
	CarID  string

	CommittedPeriods *int32 // number of hours or days booked upfront; nil for per_minute
	Status           string // "reserved" | "active" | "completed" | "cancelled"
	PricingRuleID    string
	PricingSnapshot  PricingSnapshot

	CreatedAt time.Time
	UpdatedAt time.Time
}

type BookingFilter struct {
	UserID *string
	CarID  *string

	Status        *string
	PricingRuleID *string

	Pagination *Pagination
}

type BookingCreate struct {
	UserID string
	CarID  string

	CommittedPeriods *int32
	PricingRuleID    string
}

type BookingStatusUpdate struct {
	Status string
	Reason *string
}

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

type BookingStatusReadingFilter struct {
	From       *time.Time
	To         *time.Time
	Pagination *Pagination
}
