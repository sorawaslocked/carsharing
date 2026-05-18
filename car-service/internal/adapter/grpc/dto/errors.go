package dto

import (
	"errors"

	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FromErrorToStatusCode(err error) error {
	ve, ok := errors.AsType[validation.Errors](err)
	if ok {
		return validationError(ve)
	}

	switch {
	case errors.Is(err, model.ErrMissingMetadata):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, model.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, model.ErrUnauthorized), errors.Is(err, model.ErrInsufficientPermissions):
		return status.Error(codes.PermissionDenied, err.Error())

	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, model.ErrConflict):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, model.ErrInternalServerError):
		return status.Error(codes.Internal, err.Error())

	default:
		return status.Error(codes.Internal, "something went wrong")
	}
}

func validationError(ve validation.Errors) error {
	st := status.New(codes.InvalidArgument, "invalid request")

	var fieldViolations []*errdetails.BadRequest_FieldViolation
	for field, err := range ve {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: err.Error(),
		})
	}

	st, _ = st.WithDetails(&errdetails.BadRequest{
		FieldViolations: fieldViolations,
	})

	return st.Err()
}
