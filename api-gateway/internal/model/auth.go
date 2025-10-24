package model

type Credentials struct {
	FirstName            string
	LastName             string
	Email                string
	PhoneNumber          string
	BirthDate            string
	Password             string
	PasswordConfirmation string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
