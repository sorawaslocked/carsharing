package service

import (
	"errors"

	"carsharing/user-service/internal/model"
	"github.com/go-playground/validator/v10"
)

type activationCodeValidation struct {
	Code string `validate:"required,len=6,alphanum,uppercase"`
}

func validateInput(v *validator.Validate, input any) error {
	err := v.Struct(input)
	if err == nil {
		return nil
	}

	errs := make(model.ValidationErrors)

	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	for _, fieldErr := range validationErrors {
		field := uncapitalize(fieldErr.Field())
		if _, exists := errs[field]; exists {
			continue
		}
		errs[field] = validationError(fieldErr)
	}

	return errs
}
