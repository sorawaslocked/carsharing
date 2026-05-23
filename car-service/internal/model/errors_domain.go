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
	ErrInvalidRole             = errors.New("must be a valid role")
	ErrMileageRegression       = errors.New("incoming mileage value is lower than the current mileage")
)
