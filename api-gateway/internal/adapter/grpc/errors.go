package grpc

import (
	"errors"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)

func fromGrpcErr(err error) error {
	st, ok := status.FromError(err)

	if !ok {
		return ErrInvalidStatusCode
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return ErrInvalidStatusCode
	case codes.NotFound:
		return model.ErrNotFound
	case codes.AlreadyExists:
		return model.ErrAlreadyExists
	case codes.PermissionDenied:
		return model.ErrForbidden
	case codes.Internal:
		return model.ErrInternalServerError
	case codes.Unauthenticated:
		return model.ErrUnauthorized
	default:
		return model.ErrInternalServerError
	}
}

func fromInvalidArgument(st *status.Status) error {
	return nil
}
