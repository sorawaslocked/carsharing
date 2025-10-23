package service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func validationError(fieldErr validator.FieldError) error {
	switch fieldErr.Tag() {
	case "required":
		return fmt.Errorf("required")
	case "required_without":
		return fmt.Errorf("required")
	case "max":
		return fmt.Errorf("must be at most %s characters", fieldErr.Param())
	case "min":
		return fmt.Errorf("must be at least %s characters", fieldErr.Param())
	case "email":
		return fmt.Errorf("must be a valid email address")
	case "e164":
		return fmt.Errorf("must be a valid e164 phone number")
	case "eqfield":
		return fmt.Errorf("must be same value")
	case "complex_password":
		return fmt.Errorf("must contain uppercase characters, lowercase characters, numbers, and special characters(!@#)")
	case "min_age":
		return fmt.Errorf("must be at least %s years", fieldErr.Param())
	default:
		return fmt.Errorf("validation error")
	}
}
