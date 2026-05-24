package validation

import (
	"fmt"
	"log/slog"

	sharedvalidation "carsharing/shared/validation"
	"carsharing/trip-service/internal/model"

	"github.com/go-playground/validator/v10"
)

type ErrRegisterValidator struct {
	Tag string
}

func (e ErrRegisterValidator) Error() string {
	return fmt.Sprintf("failed to register validator %q", e.Tag)
}

func RegisterCustomValidators(v *validator.Validate, log *slog.Logger) error {
	if err := sharedvalidation.RegisterLocationValidators(v, log); err != nil {
		return err
	}
	if err := sharedvalidation.RegisterTimeRangeValidators(v, log); err != nil {
		return err
	}

	validators := []struct {
		tag string
		fn  validator.Func
	}{
		{"trip_status", tripStatusValidator},
	}

	for _, vd := range validators {
		if err := v.RegisterValidation(vd.tag, vd.fn); err != nil {
			log.Error("registering validator", slog.String("tag", vd.tag), slog.Any("error", err))
			return ErrRegisterValidator{Tag: vd.tag}
		}
	}

	return nil
}

func tripStatusValidator(fl validator.FieldLevel) bool {
	switch model.TripStatus(fl.Field().String()) {
	case model.TripStatusActive, model.TripStatusCompleted, model.TripStatusCancelled:
		return true
	}
	return false
}
