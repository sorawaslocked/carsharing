package model

import "time"

type BookingStatus string

const (
	BookingStatusCreated   BookingStatus = "created"
	BookingStatusExpired   BookingStatus = "expired"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// allowedTransitions defines every legal from→to pair.
// Terminal statuses (expired, completed, cancelled) have no outgoing edges.
var allowedTransitions = map[BookingStatus][]BookingStatus{
	BookingStatusCreated: {
		BookingStatusExpired,
		BookingStatusCompleted,
		BookingStatusCancelled,
	},
}

func ValidateTransition(from, to BookingStatus) error {
	for _, allowed := range allowedTransitions[from] {
		if allowed == to {
			return nil
		}
	}
	return ErrInvalidTransition
}

func ParseBookingStatus(s string) (BookingStatus, error) {
	switch BookingStatus(s) {
	case BookingStatusCreated, BookingStatusExpired, BookingStatusCompleted, BookingStatusCancelled:
		return BookingStatus(s), nil
	default:
		return "", ErrInvalidStatus
	}
}

type Booking struct {
	ID               string
	UserID           string
	CarID            string
	CommittedPeriods *int32
	Status           BookingStatus
	PricingRuleID    string
	PricingSnapshot  PricingSnapshot
	ExpiresAt        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type BookingCreate struct {
	UserID           string
	CarID            string
	CommittedPeriods *int32
	PricingRuleID    string
}

type BookingListFilter struct {
	UserID        *string
	CarID         *string
	Status        *string
	PricingRuleID *string
	Pagination    Pagination
}
