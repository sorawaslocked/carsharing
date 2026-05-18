package model

import "time"

type Booking struct {
	ID               string
	UserID           string
	CarID            string
	Status           string
	PricingSnapshot  PricingSnapshot
	CommittedPeriods *int32
	ExpiresAt        time.Time
}
