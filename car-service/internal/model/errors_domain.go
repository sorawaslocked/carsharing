package model

import (
	"errors"
	"fmt"
)

type ErrInvalidStatusTransition struct {
	From CarStatus
	To   CarStatus
}

func (e ErrInvalidStatusTransition) Error() string {
	return fmt.Sprintf("invalid status transition: from=%s to=%s", e.From.String(), e.To.String())
}

var (
	ErrInvalidMetadata         = errors.New("invalid metadata")
	ErrUnauthenticated         = errors.New("invalid credentials")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("resource not found")
	ErrAlreadyExists           = errors.New("resource already exists")
	ErrMileageRegression       = errors.New("incoming mileage value is lower than the current mileage")

	// Entity-specific not-found errors — all wrap ErrNotFound so errors.Is(err, ErrNotFound) still holds.
	ErrCarNotFound                    = fmt.Errorf("car not found: %w", ErrNotFound)
	ErrCarModelNotFound               = fmt.Errorf("car model not found: %w", ErrNotFound)
	ErrZoneNotFound                   = fmt.Errorf("zone not found: %w", ErrNotFound)
	ErrCarInsuranceNotFound           = fmt.Errorf("car insurance not found: %w", ErrNotFound)
	ErrCarMaintenanceTemplateNotFound = fmt.Errorf("maintenance template not found: %w", ErrNotFound)
	ErrCarMaintenanceRecordNotFound   = fmt.Errorf("maintenance record not found: %w", ErrNotFound)

	// Constraint-specific duplicate errors.
	ErrDuplicateVIN          = errors.New("car with this VIN already exists")
	ErrDuplicateLicensePlate = errors.New("car with this license plate already exists")
)
