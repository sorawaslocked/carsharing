package dto

import (
	"errors"

	"carsharing/api-gateway/internal/model"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

func IsSystemErr(err error) bool {
	st, ok := status.FromError(err)
	return !ok || st.Code() == codes.Internal
}

func FromGrpcErr(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return ErrInvalidStatusCode
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return fromInvalidArgument(st)
	case codes.NotFound:
		return model.ErrNotFound
	case codes.AlreadyExists:
		return model.ErrAlreadyExists
	case codes.FailedPrecondition:
		return model.ErrConflict
	case codes.PermissionDenied:
		return model.ErrForbidden
	case codes.ResourceExhausted:
		return model.ErrTooManyRequests
	case codes.Internal:
		return model.ErrInternalServerError
	case codes.Unauthenticated:
		return model.ErrUnauthorized
	default:
		return model.ErrInternalServerError
	}
}

func fromInvalidArgument(st *status.Status) error {
	for _, d := range st.Details() {
		switch info := d.(type) {
		case *errdetails.BadRequest:
			ve := make(model.ValidationErrors)
			for _, fv := range info.FieldViolations {
				ve[fv.Field] = fv.Description
			}
			return ve
		}
	}

	return model.ErrInvalidArgument
}
