package redis

import "errors"

var (
	ErrCloseFailed  = errors.New("redis: failed to close connection")
	ErrReadFailed   = errors.New("redis: failed to read from redis")
	ErrWriteFailed  = errors.New("redis: failed to write to redis")
	ErrDeleteFailed = errors.New("redis: failed to delete from redis")
)
