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

	// Entity-specific not-found errors.
	ErrCarNotFound                    = errors.New("car not found")
	ErrCarModelNotFound               = errors.New("car model not found")
	ErrZoneNotFound                   = errors.New("zone not found")
	ErrCarInsuranceNotFound           = errors.New("car insurance not found")
	ErrCarMaintenanceTemplateNotFound = errors.New("maintenance template not found")
	ErrCarMaintenanceRecordNotFound   = errors.New("maintenance record not found")

	// Constraint-specific duplicate errors.
	ErrDuplicateVIN          = errors.New("car with this VIN already exists")
	ErrDuplicateLicensePlate = errors.New("car with this license plate already exists")
)
