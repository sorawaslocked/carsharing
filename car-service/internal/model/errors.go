package model

import "errors"

var (
	ErrInvalidMetadata = errors.New("invalid metadata")

	ErrInternalServerError     = errors.New("internal server error")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrUnauthenticated         = errors.New("unauthenticated")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("not found")
	ErrConflict                = errors.New("conflict")

	ErrInvalidRole = errors.New("must be a valid role")
)
