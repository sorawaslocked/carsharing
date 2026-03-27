package validation

import (
	"github.com/sorawaslocked/car-rental-car-service/internal/model"

	"github.com/go-playground/validator/v10"
)

const (
	maxRadiusKM float64 = 100
	minRadiusKM float64 = 0.1
)

func RegisterCustomValidators(v *validator.Validate) error {
	err := v.RegisterValidation("carfueltype", carFuelTypeValidator)
	if err != nil {
		return err
	}

	err = v.RegisterValidation("cartransmission", carTransmissionValidator)
	if err != nil {
		return err
	}

	err = v.RegisterValidation("carbodytype", carBodyTypeValidator)
	if err != nil {
		return err
	}

	err = v.RegisterValidation("carclass", carClassValidator)
	if err != nil {
		return err
	}

	err = v.RegisterValidation("carstatus", carStatusValidator)
	if err != nil {
		return err
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
