package model

import "errors"

var (
	ErrInvalidMetadata         = errors.New("invalid metadata")
	ErrUnauthenticated         = errors.New("unauthenticated")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("not found")
	ErrConflict                = errors.New("conflict")

	// Entity-specific not-found errors.
	ErrBookingNotFound     = errors.New("booking not found")
	ErrPricingRuleNotFound = errors.New("pricing rule not found")
	ErrCarNotFound         = errors.New("car not found")
	ErrCarModelNotFound    = errors.New("car model not found")
	ErrZoneNotFound        = errors.New("zone not found")

	ErrCarNotAvailable   = errors.New("car is not available for booking")
	ErrInvalidTransition = errors.New("invalid booking status transition")
	ErrInvalidStatus     = errors.New("invalid booking status")
)
