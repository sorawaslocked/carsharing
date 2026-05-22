package service

import (
	"io"
	"log/slog"
	"testing"

	"carsharing/car-service/internal/validation"
	"github.com/go-playground/validator/v10"
)

func newTestValidator(t *testing.T) *validator.Validate {
	t.Helper()
	v := validator.New()
	_ = validation.RegisterCustomValidators(v, slog.New(slog.NewTextHandler(io.Discard, nil)))
	return v
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
