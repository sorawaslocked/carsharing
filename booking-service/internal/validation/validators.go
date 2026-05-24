package validation

import (
	"fmt"
	"log/slog"

	"carsharing/booking-service/internal/model"
	sharedvalidation "carsharing/shared/validation"

	"github.com/go-playground/validator/v10"
)

type ErrRegisterValidator struct {
	Tag string
}

func (e ErrRegisterValidator) Error() string {
	return fmt.Sprintf("failed to register validator %q", e.Tag)
}

func RegisterCustomValidators(v *validator.Validate, log *slog.Logger) error {
	if err := sharedvalidation.RegisterTimeRangeValidators(v, log); err != nil {
		return err
	}

	validators := []struct {
		tag string
		fn  validator.Func
	}{
		{"booking_status", bookingStatusValidator},
		{"pricing_rule_type", pricingRuleTypeValidator},
		{"carstatus", carStatusValidator},
		{"carclass", carClassValidator},
	}

	for _, vd := range validators {
		if err := v.RegisterValidation(vd.tag, vd.fn); err != nil {
			log.Error("registering validator", slog.String("tag", vd.tag), slog.Any("error", err))
			return ErrRegisterValidator{Tag: vd.tag}
		}
	}

	return nil
}

func bookingStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseBookingStatus(fl.Field().String())
	return ok
}

func pricingRuleTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.PricingRuleTypeFromString(fl.Field().String())
	return ok
}

func carStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.CarStatusFromString(fl.Field().String())
	return ok
}

func carClassValidator(fl validator.FieldLevel) bool {
	_, ok := model.CarClassFromString(fl.Field().String())
	return ok
}
