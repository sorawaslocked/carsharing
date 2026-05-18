package minio

import "errors"

var (
	ErrFailedConnection   = errors.New("failed connection")
	ErrBucketCheckFailed  = errors.New("failed to check bucket existence")
	ErrBucketCreateFailed = errors.New("failed to create bucket")
)
