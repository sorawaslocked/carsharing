package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

var (
	ErrRequiredField = errors.New("is required")
	ErrInvalidToken  = errors.New("invalid token")
)

func validationError(fieldErr validator.FieldError) error {
	switch fieldErr.Tag() {
	case "max":
		return fmt.Errorf("must be at most %s characters", fieldErr.Param())
	case "min":
		return fmt.Errorf("must be at least %s characters", fieldErr.Param())
	case "email":
		return fmt.Errorf("must be a valid email address")
	case "e164":
		return fmt.Errorf("must be a valid e164 phone number")
	case "complex_password":
		return fmt.Errorf("must contain uppercase characters, lowercase characters, numbers, and special characters(!@#)")
	case "min_age_18":
		return fmt.Errorf("must be at least 18 years")
	default:
		return fmt.Errorf("validation error")
	}
}
