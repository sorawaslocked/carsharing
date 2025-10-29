package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

type refreshTokenValidation struct {
	RefreshToken string `validate:"required,jwt"`
}

type queryParamsValidation struct {
	ID    uint64 `validate:"required_without=Email"`
	Email string `validate:"required_without=ID"`
}

func validateInput(v *validator.Validate, input any) error {
	err := v.Struct(input)
	if err == nil {
		return nil
	}

	var errs model.ValidationErrors
	errs = make(map[string]error)
	var validationErrors validator.ValidationErrors
	errors.As(err, &validationErrors)

	for _, fieldErr := range validationErrors {
		fieldErrField := fieldErr.Field()
		field := uncapitalize(fieldErrField)

		if _, ok := errs[field]; ok {
			continue
		}
		errs[field] = validationError(fieldErr)
	}

	return errs
}

func checkQueryParams(v *validator.Validate, filter model.UserFilter) error {
	queryParams := queryParamsValidation{}

	if filter.ID != nil {
		queryParams.ID = *filter.ID
	}
	if filter.Email != nil {
		queryParams.Email = *filter.Email
	}

	return validateInput(v, queryParams)
}
