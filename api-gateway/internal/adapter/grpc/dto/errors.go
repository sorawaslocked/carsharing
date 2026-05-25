package dto

import (
	"errors"

	"carsharing/api-gateway/internal/model"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

// domainErr wraps a sentinel domain error with the original gRPC error so the
// HTTP layer can surface the service's message while errors.Is still matches.
type domainErr struct {
	sentinel error
	cause    error
}

func (e *domainErr) Error() string {
	st, ok := status.FromError(e.cause)
	if ok && st.Message() != "" {
		return st.Message()
	}
	return e.sentinel.Error()
}

func (e *domainErr) Is(target error) bool {
	return errors.Is(e.sentinel, target)
}

func (e *domainErr) Unwrap() error {
	return e.cause
}

func wrapDomain(sentinel, cause error) error {
	return &domainErr{sentinel: sentinel, cause: cause}
}

func FromGrpcErr(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return ErrInvalidStatusCode
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return fromInvalidArgument(st, err)
	case codes.NotFound:
		return wrapDomain(model.ErrNotFound, err)
	case codes.AlreadyExists:
		return wrapDomain(model.ErrAlreadyExists, err)
	case codes.FailedPrecondition, codes.Aborted:
		return wrapDomain(model.ErrConflict, err)
	case codes.PermissionDenied:
		return wrapDomain(model.ErrForbidden, err)
	case codes.ResourceExhausted:
		return wrapDomain(model.ErrTooManyRequests, err)
	case codes.Internal:
		return wrapDomain(model.ErrInternalServerError, err)
	case codes.Unauthenticated:
		return wrapDomain(model.ErrUnauthorized, err)
	default:
		return wrapDomain(model.ErrInternalServerError, err)
	}
}

func fromInvalidArgument(st *status.Status, err error) error {
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

	return wrapDomain(model.ErrInvalidArgument, err)
}
