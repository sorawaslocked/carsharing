package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"strings"
	"time"
)

type registerValidation struct {
	Email                string    `validate:"required,email"`
	PhoneNumber          string    `validate:"required,e164"`
	Password             string    `validate:"required,min=8,max=20,complex_password"`
	PasswordConfirmation string    `validate:"required,min=8,max=20,eqfield=Password,complex_password"`
	FirstName            string    `validate:"required,min=1,max=100"`
	LastName             string    `validate:"required,min=1,max=100"`
	BirthDate            time.Time `validate:"required,min_age=18"`
}

type loginValidation struct {
	Email       string `validate:"required_without=PhoneNumber,omitempty,email"`
	PhoneNumber string `validate:"required_without=Email,omitempty,e164"`
	Password    string `validate:"required"`
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
		field := strings.ToLower(fieldErrField[:1]) + fieldErrField[1:]

		if _, ok := errs[field]; ok {
			continue
		}
		errs[field] = validationError(fieldErr)
	}

	return errs
}

func toRoleStrings(roles []model.Role) []string {
	var result []string
	for _, role := range roles {
		result = append(result, role.String())
	}

	return result
}
