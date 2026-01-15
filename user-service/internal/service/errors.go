package service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

func validationError(fieldErr validator.FieldError) error {
	switch fieldErr.Tag() {
	case "required":
		return model.ErrRequiredField
	case "required_without":
		field := uncapitalize(fieldErr.Field())
		param := uncapitalize(fieldErr.Param())

		return fmt.Errorf("either %s is required or %s", field, param)
	case "required_with":
		field := uncapitalize(fieldErr.Field())
		param := uncapitalize(fieldErr.Param())

		return fmt.Errorf("%s required with %s", field, param)
	case "eqfield":
		param := uncapitalize(fieldErr.Param())

		return fmt.Errorf("must be same value as %s", param)
	case "alpha":
		return model.ErrNotAlpha
	case "alphaunicode":
		return model.ErrNotAlphaUnicode
	case "uppercase":
		return model.ErrNotUppercase
	case "len":
		return fmt.Errorf("must be exactly %s characters long", fieldErr.Field())
	case "max":
		return fmt.Errorf("must be at most %s characters", fieldErr.Param())
	case "min":
		return fmt.Errorf("must be at least %s characters", fieldErr.Param())
	case "email":
		return model.ErrInvalidEmail
	case "e164":
		return model.ErrInvalidPhoneNumber
	case "jwt":
		return model.ErrInvalidJwtToken
	case "complex_password":
		return model.ErrNotComplexPassword
	case "min_age":
		return fmt.Errorf("must be at least %s years", fieldErr.Param())
	default:
		return fmt.Errorf("validation error")
	}
}
