package validation

import (
	"errors"
	"fmt"

	"carsharing/shared/pkg/utils"

	"github.com/go-playground/validator/v10"
)

func ValidateInput(v *validator.Validate, input any) error {
	err := v.Struct(input)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	errs := make(Errors)
	for _, fieldErr := range validationErrors {
		field := utils.Uncapitalize(fieldErr.Field())
		if _, exists := errs[field]; exists {
			continue
		}
		errs[field] = validationError(fieldErr)
	}

	return errs
}

func ValidateActivationCode(v *validator.Validate, code string) error {
	return ValidateInput(v, activationCodeValidation{Code: code})
}

func ValidateID(v *validator.Validate, id string) error {
	return ValidateInput(v, idValidation{ID: id})
}

func validationError(fieldErr validator.FieldError) error {
	switch fieldErr.Tag() {
	case "required":
		return ErrRequiredField
	case "required_without":
		field := utils.Uncapitalize(fieldErr.Field())
		param := utils.Uncapitalize(fieldErr.Param())
		return fmt.Errorf("either %s or %s is required", field, param)
	case "required_with":
		field := utils.Uncapitalize(fieldErr.Field())
		param := utils.Uncapitalize(fieldErr.Param())
		return fmt.Errorf("%s is required when %s is set", field, param)
	case "eqfield":
		param := utils.Uncapitalize(fieldErr.Param())
		return fmt.Errorf("must equal %s", param)
	case "alphanum":
		return ErrNotAlphaNum
	case "alphaunicode":
		return ErrNotAlphaUnicode
	case "uppercase":
		return ErrNotUppercase
	case "len":
		return fmt.Errorf("must be exactly %s characters", fieldErr.Param())
	case "max":
		return fmt.Errorf("must be at most %s characters", fieldErr.Param())
	case "min":
		return fmt.Errorf("must be at least %s characters", fieldErr.Param())
	case "email":
		return ErrInvalidEmail
	case "e164":
		return ErrInvalidPhoneNumber
	case "complex_password":
		return ErrNotComplexPassword
	case "min_age":
		return fmt.Errorf("must be at least %s years old", fieldErr.Param())
	case "document_image_type":
		return ErrInvalidImageType
	case "document_status":
		return ErrInvalidDocumentStatus
	case "document_sort":
		return ErrInvalidDocumentSort
	case "uuid4":
		return ErrInvalidID
	case "role":
		return ErrInvalidRole
	default:
		return fmt.Errorf("invalid value")
	}
}
