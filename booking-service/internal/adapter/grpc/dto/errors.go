package dto

import (
	"errors"

	"github.com/sorawaslocked/car-rental-booking-service/internal/model"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, model.ErrConflict):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, model.ErrInsufficientPermissions):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, model.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, model.ErrInvalidTransition):
		return fieldViolation("status", err.Error())
	case errors.Is(err, model.ErrInvalidStatus):
		return fieldViolation("status", err.Error())
	default:
		return status.Error(codes.Internal, model.ErrInternalServerError.Error())
	}
}

func fieldViolation(field, desc string) error {
	st, _ := status.New(codes.InvalidArgument, "validation failed").
		WithDetails(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{Field: field, Description: desc},
			},
		})
	return st.Err()
}
