package grpc

import (
	"errors"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
		return fromInvalidArgument(st)
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
	for _, d := range st.Details() {
		switch info := d.(type) {
		case *errdetails.BadRequest:
			var ve model.ValidationErrors
			ve = make(map[string]string)

			for _, fv := range info.FieldViolations {
				ve[fv.Field] = fv.Description
			}

			return ve
		}
	}

	return model.ErrInvalidArgument
}
