package model

import "errors"

var (
	ErrMissingMetadata = errors.New("missing metadata")

	ErrInternalServerError = errors.New("internal server error")
	ErrUnauthorized        = errors.New("unauthorized")
)
