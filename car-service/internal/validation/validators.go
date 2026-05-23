package validation

import (
	"fmt"
	"log/slog"

	"carsharing/car-service/internal/model"

	"github.com/go-playground/validator/v10"
)

const (
	maxRadiusKM float64 = 100
	minRadiusKM float64 = 0.1
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
		{"carfueltype", carFuelTypeValidator},
		{"cartransmission", carTransmissionValidator},
		{"carbodytype", carBodyTypeValidator},
		{"carclass", carClassValidator},
		{"carstatus", carStatusValidator},
		{"zonetype", zoneTypeValidator},
		{"insurancetype", insuranceTypeValidator},
		{"insurancestatus", insuranceStatusValidator},
		{"maintenancerecordstatus", maintenanceRecordStatusValidator},
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

func carFuelTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.CarFuelTypeFromString(fl.Field().String())
	return ok
}

func carTransmissionValidator(fl validator.FieldLevel) bool {
	_, ok := model.CarTransmissionFromString(fl.Field().String())
	return ok
}

func carBodyTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.CarBodyTypeFromString(fl.Field().String())
	return ok
}

func carClassValidator(fl validator.FieldLevel) bool {
	_, ok := model.CarClassFromString(fl.Field().String())
	return ok
}

func carStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.CarStatusFromString(fl.Field().String())
	return ok
}

func zoneTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.ZoneTypeFromString(fl.Field().String())
	return ok
}

func insuranceTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.InsuranceTypeFromString(fl.Field().String())
	return ok
}

func insuranceStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.InsuranceStatusFromString(fl.Field().String())
	return ok
}

func maintenanceRecordStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.MaintenanceRecordStatusFromString(fl.Field().String())
	return ok
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
	return r >= minRadiusKM && r <= maxRadiusKM
}
