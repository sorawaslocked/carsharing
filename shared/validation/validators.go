package validation

import (
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
)

const (
	MaxRadiusKM float64 = 100
	MinRadiusKM float64 = 0.1
)

type ErrRegisterValidator struct {
	Tag string
}

func (e ErrRegisterValidator) Error() string {
	return fmt.Sprintf("failed to register validator %q", e.Tag)
}

func RegisterLocationValidators(v *validator.Validate, log *slog.Logger) error {
	validators := []struct {
		tag string
		fn  validator.Func
	}{
		{"latitude_range", latitudeValidator},
		{"longitude_range", longitudeValidator},
		{"radius_range", radiusValidator},
	}

	for _, vd := range validators {
		if err := v.RegisterValidation(vd.tag, vd.fn); err != nil {
			log.Error("registering validator", slog.String("tag", vd.tag), slog.Any("error", err))
			return ErrRegisterValidator{Tag: vd.tag}
		}
	}

	return nil
}

func RegisterTimeRangeValidators(v *validator.Validate, _ *slog.Logger) error {
	v.RegisterStructValidation(validateTimeRange, TimeRange{})
	return nil
}

func validateTimeRange(sl validator.StructLevel) {
	tr := sl.Current().Interface().(TimeRange)
	if tr.From != nil && tr.To != nil && tr.From.After(*tr.To) {
		sl.ReportError(tr.To, "To", "To", "time_range_from_before_to", "")
	}
}

func latitudeValidator(fl validator.FieldLevel) bool {
	lat := fl.Field().Float()
	return lat >= -90 && lat <= 90
}

func longitudeValidator(fl validator.FieldLevel) bool {
	lon := fl.Field().Float()
	return lon >= -180 && lon <= 180
}

func radiusValidator(fl validator.FieldLevel) bool {
	r := fl.Field().Float()
	return r >= MinRadiusKM && r <= MaxRadiusKM
}
