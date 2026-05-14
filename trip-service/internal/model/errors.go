package model

import "errors"

var (
	ErrUnauthenticated         = errors.New("unauthenticated")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("not found")
	ErrAlreadyExists           = errors.New("already exists")
	ErrBookingNotReserved      = errors.New("booking is not in reserved status")
	ErrTripNotActive           = errors.New("trip is not active")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
