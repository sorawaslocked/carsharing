package dto_test

import (
	"errors"
	"testing"

	"carsharing/car-service/internal/adapter/grpc/dto"
	"carsharing/car-service/internal/model"
	"carsharing/car-service/internal/validation"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFromErrorToStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode codes.Code
	}{
		{"not found", model.ErrNotFound, codes.NotFound},
		{"conflict", model.ErrConflict, codes.AlreadyExists},
		{"internal server error", model.ErrInternalServerError, codes.Internal},
		{"missing metadata", model.ErrMissingMetadata, codes.InvalidArgument},
		{"unauthenticated", model.ErrUnauthenticated, codes.Unauthenticated},
		{"unauthorized", model.ErrUnauthorized, codes.PermissionDenied},
		{"insufficient permissions", model.ErrInsufficientPermissions, codes.PermissionDenied},
		{"unknown error defaults to internal", errors.New("something unexpected"), codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dto.FromErrorToStatusCode(tt.err)
			assert.Equal(t, tt.wantCode, status.Code(got))
		})
	}

	t.Run("validation errors map to InvalidArgument with field details", func(t *testing.T) {
		ve := make(validation.Errors)
		ve["name"] = errors.New("this field is required")

		got := dto.FromErrorToStatusCode(ve)
		assert.Equal(t, codes.InvalidArgument, status.Code(got))

		st, ok := status.FromError(got)
		assert.True(t, ok)
		assert.NotEmpty(t, st.Details())
	})
}
