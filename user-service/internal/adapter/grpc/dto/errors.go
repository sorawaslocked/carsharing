package dto

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/validation"
)

func validationError(ve validation.Errors) error {
	st := status.New(codes.InvalidArgument, "validation failed")

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

func ToStatusError(err error) error {
	var ve validation.Errors

	switch {
	case errors.Is(err, model.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, model.ErrInsufficientPermissions):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, model.ErrDuplicateEmail),
		errors.Is(err, model.ErrDuplicatePhone),
		errors.Is(err, model.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, model.ErrNoUpdateFields):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, model.ErrActivationCodeResendTooSoon):
		return status.Error(codes.ResourceExhausted, err.Error())
	case errors.Is(err, validation.ErrInvalidEmail):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.As(err, &ve):
		return validationError(ve)
	default:
		return status.Error(codes.Internal, "something went wrong")
	}
}
