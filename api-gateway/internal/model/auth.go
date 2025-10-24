package model

type Credentials struct {
	Email                string
	PhoneNumber          string
	Password             string
	PasswordConfirmation string
	FirstName            string
	LastName             string
	BirthDate            string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
