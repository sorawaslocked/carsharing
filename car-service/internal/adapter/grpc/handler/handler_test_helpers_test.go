package handler

import (
	"errors"
	"io"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// errInternal is an unrecognized error that maps to codes.Internal via the default case.
var errInternal = errors.New("unexpected internal error")

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func grpcCode(err error) codes.Code {
	return status.Code(err)
}
