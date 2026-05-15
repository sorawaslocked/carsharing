package service

import (
	"io"
	"log/slog"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
	"github.com/sorawaslocked/car-rental-car-service/internal/validation"
)

// newTestValidator builds a validator with all custom tags registered,
// including zone, insurance, and maintenance tags missing from the production
// RegisterCustomValidators (which only registers car-specific tags).
func newTestValidator(t *testing.T) *validator.Validate {
	t.Helper()
	v := validator.New()
	_ = validation.RegisterCustomValidators(v)
	_ = v.RegisterValidation("zonetype", func(fl validator.FieldLevel) bool {
		_, ok := model.ParseZoneType(fl.Field().String())
		return ok
	})
	_ = v.RegisterValidation("insurancetype", func(fl validator.FieldLevel) bool {
		_, ok := model.ParseInsuranceType(fl.Field().String())
		return ok
	})
	_ = v.RegisterValidation("insurancestatus", func(fl validator.FieldLevel) bool {
		_, ok := model.ParseInsuranceStatus(fl.Field().String())
		return ok
	})
	_ = v.RegisterValidation("maintenancerecordstatus", func(fl validator.FieldLevel) bool {
		_, ok := model.ParseMaintenanceRecordStatus(fl.Field().String())
		return ok
	})
	return v
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
