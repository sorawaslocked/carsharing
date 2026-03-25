package validation

import (
	"car-rental-car-service/internal/pkg/utils"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateInput(v *validator.Validate, input any) error {
	err := v.Struct(input)

	ve, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		return err
	}

	var errs Errors
	errs = make(map[string]error)

	for _, fe := range ve {
		field := utils.Uncapitalize(fe.Field())

		if _, ok := errs[field]; ok {
			continue
		}

		errs[field] = validationError(fe)
	}

	return errs
}

func validationError(fe validator.FieldError) error {
	switch fe.Tag() {
	// General
	case "required":
		return ErrRequiredField
	case "min":
		return fmt.Errorf("must be at least %s", fe.Param())
	case "max":
		return fmt.Errorf("must be at most %s", fe.Param())
	// Car
	case "carfueltype":
		return ErrInvalidCarFuelType
	case "cartransmission":
		return ErrInvalidCarTransmission
	case "carbodytype":
		return ErrInvalidCarBodyType
	case "carclass":
		return ErrInvalidCarClass
	case "carstatus":
		return ErrInvalidCarStatus
	// Location
	case "latitude_range":
		return ErrInvalidLatitudeRange
	case "longitude_range":
		return ErrInvalidLongitudeRange
	case "radius_range":
		return ErrInvalidRadiusRange
	default:
		return fmt.Errorf("failed validation: %s", fe.Tag())
	}
}
