package dto

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"carsharing/trip-service/internal/model"
)

func ToStatusError(err error) error {
	var ve model.ValidationErrors
	if errors.As(err, &ve) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	switch {
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, model.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, model.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, model.ErrInsufficientPermissions):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, model.ErrBookingNotReserved),
		errors.Is(err, model.ErrInvalidStatusTransition),
		errors.Is(err, model.ErrTripNotActive):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
