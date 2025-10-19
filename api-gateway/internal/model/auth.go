package model

type Credentials struct {
	Email                *string
	PhoneNumber          *string
	Password             string
	PasswordConfirmation *string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
