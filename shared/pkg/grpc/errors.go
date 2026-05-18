package grpc

import "errors"

var (
	ErrFailedConnection = errors.New("failed connection")
	ErrInvalidTLSFiles  = errors.New("invalid tls files")
)
