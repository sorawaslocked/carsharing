package dto

import (
	"errors"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGrpcValidationError(validationErrors model.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for field, err := range validationErrors {
		errs[field] = err.Error()
	}

	return errs
}

func ToStatusCodeError(err error) error {
	var ve model.ValidationErrors

	switch {
	case errors.As(err, &ve):
		return status.Error(codes.InvalidArgument, "invalid request")
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, "resource not found")
	default:
		return status.Error(codes.Internal, "something went wrong")
	}
}
