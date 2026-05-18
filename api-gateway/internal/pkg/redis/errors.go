package redis

import "errors"

var (
	ErrFailedConnection = errors.New("redis: failed to connect")
)
