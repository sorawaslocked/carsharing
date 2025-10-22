package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtProvider struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJwtProvider(
	secretKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *JwtProvider {
	return &JwtProvider{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type Claims struct {
	UserID uint64
	Role   string
}

func (jp *JwtProvider) GenerateAccessToken(userID uint64, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(jp.accessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jp.secretKey))
}

func (jp *JwtProvider) GenerateRefreshToken(userID uint64, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(jp.refreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jp.secretKey))
}

func (jp *JwtProvider) VerifyAndParseClaims(token string) (Claims, error) {
	jwtClaims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jp.secretKey), nil
	})
	if err != nil {
		return Claims{}, err
	}

	claims := Claims{}
	claims.UserID = uint64(jwtClaims["user_id"].(float64))
	claims.Role = jwtClaims["role"].(string)

	return claims, nil
}
