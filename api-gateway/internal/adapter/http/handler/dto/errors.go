package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

func errorBody(message string, metadata map[string]any) map[string]any {
	body := map[string]any{
		"error": map[string]any{
			"message": message,
		},
	}

	for k, v := range metadata {
		body["error"].(map[string]any)[k] = v
	}

	return body
}

func FromError(ctx *gin.Context, err error) {
	var ve model.ValidationErrors

	switch {
	case errors.As(err, &ve):
		validationError(ctx, ve)
	case errors.Is(err, model.ErrUnauthorized):
		unauthorized(ctx)
	case errors.Is(err, model.ErrForbidden):
		forbidden(ctx)
	case errors.Is(err, model.ErrNotFound):
		notFound(ctx)
	case errors.Is(err, model.ErrAlreadyExists):
		conflict(ctx)
	case errors.Is(err, model.ErrInternalServerError):
		internalServerError(ctx)
	default:
		internalServerError(ctx)
	}
}
