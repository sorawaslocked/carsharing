package model

import "time"

type Credentials struct {
	Email                string
	PhoneNumber          string
	Password             string
	PasswordConfirmation string
	FirstName            string
	LastName             string
	BirthDate            time.Time
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
