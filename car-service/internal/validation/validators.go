package validation

import (
	"log/slog"

	"carsharing/car-service/internal/model"
	sharedvalidation "carsharing/shared/validation"

	"github.com/go-playground/validator/v10"
)

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
		{"carfueltype", carFuelTypeValidator},
		{"cartransmission", carTransmissionValidator},
		{"carbodytype", carBodyTypeValidator},
		{"carclass", carClassValidator},
		{"carstatus", carStatusValidator},
		{"zonetype", zoneTypeValidator},
		{"insurancetype", insuranceTypeValidator},
		{"insurancestatus", insuranceStatusValidator},
		{"maintenancerecordstatus", maintenanceRecordStatusValidator},
	}

	for _, vd := range validators {
		if err := v.RegisterValidation(vd.tag, vd.fn); err != nil {
			log.Error("registering validator", slog.String("tag", vd.tag), slog.Any("error", err))
			return sharedvalidation.ErrRegisterValidator{Tag: vd.tag}
		}
	}

	v.RegisterStructValidation(carInsuranceUpdateValidator, CarInsuranceUpdate{})

	return nil
}

func carInsuranceUpdateValidator(sl validator.StructLevel) {
	u := sl.Current().Interface().(CarInsuranceUpdate)
	if u.StartsAt != nil && u.ExpiresAt != nil && !u.ExpiresAt.After(*u.StartsAt) {
		sl.ReportError(u.ExpiresAt, "ExpiresAt", "ExpiresAt", "gtfield", "StartsAt")
	}
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
