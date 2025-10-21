package model

type Credentials struct {
	Email                *string
	PhoneNumber          *string
	Password             string
	PasswordConfirmation *string
	FirstName            *string
	LastName             *string
	DateOfBirth          *string
}
