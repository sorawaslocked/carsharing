package model

type Credentials struct {
	FirstName            *string
	LastName             *string
	DateOfBirth          *string
	Email                *string
	PhoneNumber          *string
	Password             string
	PasswordConfirmation *string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
