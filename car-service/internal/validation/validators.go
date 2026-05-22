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
	}

	for _, vd := range validators {
		if err := v.RegisterValidation(vd.tag, vd.fn); err != nil {
			log.Error("registering validator", slog.String("tag", vd.tag), slog.Any("error", err))
			return ErrRegisterValidator{Tag: vd.tag}
		}
	}

	v.RegisterStructValidation(locationFilterValidator, model.LocationFilter{})

	return nil
}

func carFuelTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseCarFuelType(fl.Field().String())

	return ok
}

func carTransmissionValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseCarTransmission(fl.Field().String())

	return ok
}

func carBodyTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseCarBodyType(fl.Field().String())

	return ok
}

func carClassValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseCarClass(fl.Field().String())

	return ok
}

func carStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseCarStatus(fl.Field().String())

	return ok
}

func zoneTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseZoneType(fl.Field().String())
	return ok
}

func insuranceTypeValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseInsuranceType(fl.Field().String())
	return ok
}

func insuranceStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseInsuranceStatus(fl.Field().String())
	return ok
}

func maintenanceRecordStatusValidator(fl validator.FieldLevel) bool {
	_, ok := model.ParseMaintenanceRecordStatus(fl.Field().String())
	return ok
}

func locationFilterValidator(sl validator.StructLevel) {
	lf := sl.Current().Interface().(model.LocationFilter)

	if lf.Location.Latitude < -90 || lf.Location.Latitude > 90 {
		sl.ReportError(lf.Location, "Location", "latitude", "latitude_range", "")
	}
	if lf.Location.Longitude < -180 || lf.Location.Longitude > 180 {
		sl.ReportError(lf.Location.Longitude, "Location.Longitude", "longitude", "longitude_range", "")
	}
	if lf.RadiusKM < minRadiusKM || lf.RadiusKM > maxRadiusKM {
		sl.ReportError(lf.RadiusKM, "RadiusKM", "radiuskm", "radius_range", "")
	}
}
