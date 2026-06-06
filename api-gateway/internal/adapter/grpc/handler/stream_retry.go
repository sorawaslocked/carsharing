package handler

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const streamReconnectDelay = 2 * time.Second

func isUnavailable(err error) bool {
	s, ok := status.FromError(err)
	return ok && s.Code() == codes.Unavailable
}
