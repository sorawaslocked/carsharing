package model

import "errors"

var (
	ErrUnauthenticated         = errors.New("invalid credentials")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrNotFound                = errors.New("resource not found")
	ErrNoUpdateFields          = errors.New("no update fields provided")
	ErrAlreadyExists           = errors.New("resource already exists")
	ErrDuplicateEmail          = errors.New("user with this email already exists")
	ErrDuplicatePhone          = errors.New("user with this phone number already exists")
)
