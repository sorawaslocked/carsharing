package postgres

import "errors"

var (
	ErrNoUpdateFields = errors.New("no update fields set")
)
