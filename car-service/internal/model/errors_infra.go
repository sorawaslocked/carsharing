package model

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrObjectStorage       = errors.New("object storage error")
)
