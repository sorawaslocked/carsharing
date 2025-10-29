package dto

import (
	"errors"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validationError(ve model.ValidationErrors) error {
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

func ToStatusCodeError(err error) error {
	var ve model.ValidationErrors

	switch {
	case errors.As(err, &ve):
		return validationError(ve)
	case errors.Is(err, model.ErrNoUpdateFields):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, "resource not found")
	case errors.Is(err, model.ErrDuplicateEmail):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "something went wrong")
	}
}
