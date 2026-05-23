package validation

import (
	"fmt"
	"log/slog"

	"carsharing/booking-service/internal/model"

	"github.com/go-playground/validator/v10"
)

type ErrRegisterValidator struct {
	Tag string
}

func (e ErrRegisterValidator) Error() string {
	return fmt.Sprintf("failed to register validator %q", e.Tag)
}

func RegisterCustomValidators(v *validator.Validate, log *slog.Logger) error {
	validators := []struct {
		tag string
		fn  validator.Func
	}{
		{"booking_status", bookingStatusValidator},
		{"pricing_rule_type", pricingRuleTypeValidator},
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
	_, err := model.ParseBookingStatus(fl.Field().String())
	return err == nil
}

func pricingRuleTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.PricingRuleTypeFromString(fl.Field().String())
	return ok
}
