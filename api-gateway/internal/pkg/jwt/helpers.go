package jwt

import "github.com/golang-jwt/jwt/v5"

type customClaims struct {
	DeviceID string `json:"device_id"`
	jwt.RegisteredClaims
}
