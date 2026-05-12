package model

import "errors"

var (
	ErrMissingMetadata         = errors.New("missing metadata")
	ErrUnauthenticated         = errors.New("invalid credentials")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("resource not found")
	ErrNoUpdateFields          = errors.New("no update fields provided")
	ErrAlreadyExists           = errors.New("resource already exists")
)
