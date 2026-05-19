package model

import "errors"

var (
	ErrUnauthenticated         = errors.New("unauthenticated")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("not found")
	ErrAlreadyExists           = errors.New("already exists")
	ErrBookingNotCreated       = errors.New("booking is not in created status")
	ErrTripNotActive           = errors.New("trip is not active")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
