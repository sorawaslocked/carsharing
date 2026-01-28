package model

type Credentials struct {
	Email       string `validate:"required_without=PhoneNumber,omitempty,email"`
	PhoneNumber string `validate:"required_without=Email,omitempty,e164"`
	Password    string `validate:"required"`
}

type Token struct {
	AccessToken           string
	AccessTokenExpiresIn  int64
	RefreshToken          string
	RefreshTokenExpiresIn int64
}
