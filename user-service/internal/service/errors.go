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

		return fmt.Errorf("%s required without %s", field, param)
	case "eqfield":
		param := uncapitalize(fieldErr.Param())

		return fmt.Errorf("must be same value as %s", param)
	case "alphaunicode":
		return model.ErrNotAlphaUnicode
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
		return fmt.Errorf("must contain uppercase characters, lowercase characters, numbers, and special characters(!@#)")
	case "min_age":
		return fmt.Errorf("must be at least %s years", fieldErr.Param())
	default:
		return fmt.Errorf("validation error")
	}
}
