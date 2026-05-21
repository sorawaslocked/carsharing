package model

import "errors"

var (
	ErrInvalidMetadata             = errors.New("invalid metadata")
	ErrUnauthenticated             = errors.New("invalid credentials")
	ErrInsufficientPermissions     = errors.New("insufficient permissions")
	ErrNotFound                    = errors.New("resource not found")
	ErrNoUpdateFields              = errors.New("no update fields provided")
	ErrAlreadyExists               = errors.New("resource already exists")
	ErrDuplicateEmail              = errors.New("user with this email already exists")
	ErrDuplicatePhone              = errors.New("user with this phone number already exists")
	ErrEmailVerified               = errors.New("user with this email already verified")
	ErrActivationCodeResendTooSoon = errors.New("activation code was already sent recently, please wait before requesting a new one")
)
