package validation

import (
	"errors"
	"fmt"

	"carsharing/shared/pkg/utils"
	sharedvalidation "carsharing/shared/validation"

	"github.com/go-playground/validator/v10"
)

type idValidation struct {
	ID string `validate:"required,uuid4"`
}

func ValidateID(v *validator.Validate, id string) error {
	return ValidateInput(v, idValidation{ID: id})
}

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
	case "uuid4", "uuid":
		return ErrInvalidID
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
		return sharedvalidation.ErrInvalidLatitudeRange
	case "longitude_range":
		return sharedvalidation.ErrInvalidLongitudeRange
	case "radius_range":
		return sharedvalidation.ErrInvalidRadiusRange
	default:
		return fmt.Errorf("failed validation: %s", fe.Tag())
	}
}
