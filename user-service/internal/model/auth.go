package model

import "time"

type Credentials struct {
	Email                *string    `validate:"email"`
	PhoneNumber          *string    `validate:"e164"`
	Password             string     `validate:"min=8,max=20,complex_password"`
	PasswordConfirmation *string    `validate:"min=8,max=20,complex_password"`
	FirstName            *string    `validate:"min=1,max=100"`
	LastName             *string    `validate:"min=1,max=100"`
	BirthDate            *time.Time `validate:"min_age_18"`
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
