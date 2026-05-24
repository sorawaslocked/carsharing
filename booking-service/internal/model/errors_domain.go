package model

import "errors"

var (
	ErrInvalidMetadata         = errors.New("invalid metadata")
	ErrUnauthenticated         = errors.New("unauthenticated")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("not found")
	ErrConflict                = errors.New("conflict")

	ErrInvalidTransition = errors.New("invalid booking status transition")
	ErrInvalidStatus     = errors.New("invalid booking status")
)
