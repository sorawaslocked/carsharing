package service

import (
	"fmt"

	"carsharing/user-service/internal/model"
	"github.com/go-playground/validator/v10"
)

func validationError(fieldErr validator.FieldError) error {
	switch fieldErr.Tag() {
	case "required":
		return model.ErrRequiredField
	case "required_without":
		field := uncapitalize(fieldErr.Field())
		param := uncapitalize(fieldErr.Param())
		return fmt.Errorf("either %s or %s is required", field, param)
	case "required_with":
		field := uncapitalize(fieldErr.Field())
		param := uncapitalize(fieldErr.Param())
		return fmt.Errorf("%s is required when %s is set", field, param)
	case "eqfield":
		param := uncapitalize(fieldErr.Param())
		return fmt.Errorf("must equal %s", param)
	case "alphanum":
		return model.ErrNotAlphaNum
	case "alphaunicode":
		return model.ErrNotAlphaUnicode
	case "uppercase":
		return model.ErrNotUppercase
	case "len":
		return fmt.Errorf("must be exactly %s characters", fieldErr.Param())
	case "max":
		return fmt.Errorf("must be at most %s characters", fieldErr.Param())
	case "min":
		return fmt.Errorf("must be at least %s characters", fieldErr.Param())
	case "email":
		return model.ErrInvalidEmail
	case "e164":
		return model.ErrInvalidPhoneNumber
	case "complex_password":
		return model.ErrNotComplexPassword
	case "min_age":
		return fmt.Errorf("must be at least %s years old", fieldErr.Param())
	default:
		return fmt.Errorf("invalid value")
	}
}
