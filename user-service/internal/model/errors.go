package model

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrRequiredField       = errors.New("is required")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
	ErrInvalidToken        = errors.New("invalid token")
	ErrJwt                 = errors.New("jwt error")
	ErrBcrypt              = errors.New("bcrypt error")
)
