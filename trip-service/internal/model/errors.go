package model

import "errors"

var (
	ErrInvalidMetadata             = errors.New("invalid metadata")
	ErrUnauthenticated             = errors.New("unauthenticated")
	ErrInsufficientPermissions     = errors.New("insufficient permissions")
	ErrNotFound                    = errors.New("not found")
	ErrAlreadyExists               = errors.New("already exists")
	ErrConflict                    = errors.New("conflict")
	ErrBookingNotCreated           = errors.New("booking is not in created status")
	ErrLocationInNoDropZone        = errors.New("location is within a no-drop zone")
	ErrTripNotActive               = errors.New("trip is not active")
	ErrTripNotCompleted            = errors.New("trip is not completed")
	ErrInvalidTripStatusTransition = errors.New("invalid status transition")
)
