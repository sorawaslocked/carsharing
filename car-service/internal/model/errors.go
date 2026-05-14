package model

import "errors"

var (
	ErrMissingMetadata = errors.New("missing metadata")

	ErrInternalServerError     = errors.New("internal server error")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrUnauthenticated         = errors.New("unauthenticated")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("not found")
	ErrConflict                = errors.New("conflict")

	ErrInvalidRole = errors.New("must be a valid role")
)
