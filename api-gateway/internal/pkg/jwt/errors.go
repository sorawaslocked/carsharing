package jwt

import "errors"

var (
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrInvalidToken          = errors.New("invalid token")
	ErrExpiredToken          = errors.New("expired token")
)
