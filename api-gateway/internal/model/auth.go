package model

type Credentials struct {
	Email       string
	PhoneNumber string
	Password    string
}

type Token struct {
	AccessToken  string
	RefreshToken string
}
