package grpc

import "errors"

var (
	ErrGrpcServerOffline = errors.New("grpc server is offline")
)
