package dto

type RegisterRequest struct {
	Email                string `json:"email"`
	PhoneNumber          string `json:"phoneNumber"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	DateOfBirth          string `json:"dateOfBirth"`
}

type RegisterResponse struct {
	ID     *uint64           `json:"id,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

type LoginRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
	Password    string  `json:"password"`
}

type LoginResponse struct {
	AccessToken  *string           `json:"accessToken,omitempty"`
	RefreshToken *string           `json:"refreshToken,omitempty"`
	Errors       map[string]string `json:"errors,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken  *string           `json:"accessToken,omitempty"`
	RefreshToken *string           `json:"refreshToken,omitempty"`
	Errors       map[string]string `json:"errors,omitempty"`
}
