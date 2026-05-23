package validation

import (
	"time"

	sharedvalidation "carsharing/shared/validation"
)

type UserCreate struct {
	Email                string    `validate:"required,email"`
	PhoneNumber          *string   `validate:"omitempty,e164"`
	FirstName            string    `validate:"required,min=1,max=100"`
	LastName             string    `validate:"required,min=1,max=100"`
	BirthDate            time.Time `validate:"required,min_age=18"`
	Password             string    `validate:"required,min=8,max=20,complex_password"`
	PasswordConfirmation string    `validate:"required,min=8,max=20,complex_password,eqfield=Password"`
}

type UserUpdate struct {
	Email                *string    `validate:"omitempty,email"`
	PhoneNumber          *string    `validate:"omitempty,e164"`
	FirstName            *string    `validate:"omitempty,min=1,max=100"`
	LastName             *string    `validate:"omitempty,min=1,max=100"`
	BirthDate            *time.Time `validate:"omitempty,min_age=18"`
	Password             *string    `validate:"omitempty,min=8,max=20,complex_password"`
	PasswordConfirmation *string    `validate:"omitempty,required_with=Password,min=8,max=20,complex_password,eqfield=Password"`
	ProfileImageKey      *string    `validate:"omitempty,min=1"`

	Roles              []string `validate:"omitempty,min=1,dive,role"`
	IsDocumentVerified *bool
	IsEmailVerified    *bool
	IsSuspended        *bool
}

type UserFilter struct {
	Email              *string `validate:"omitempty,email"`
	PhoneNumber        *string `validate:"omitempty,e164"`
	FirstName          *string `validate:"omitempty,min=1,max=100"`
	LastName           *string `validate:"omitempty,min=1,max=100"`
	IsDocumentVerified *bool
	IsEmailVerified    *bool
	IsSuspended        *bool
	Pagination         *sharedvalidation.Pagination
}

type Credentials struct {
	Email       *string `validate:"required_without=PhoneNumber,omitempty,email"`
	PhoneNumber *string `validate:"required_without=Email,omitempty,e164"`
	Password    string  `validate:"required"`
}

type activationCodeValidation struct {
	Code string `validate:"required,len=6,alphanum,uppercase"`
}

type idValidation struct {
	ID string `validate:"required,uuid4"`
}
