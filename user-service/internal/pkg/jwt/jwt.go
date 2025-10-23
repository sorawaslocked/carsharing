package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Provider struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewProvider(
	secretKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Provider {
	return &Provider{
		secretKey:       secretKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type Claims struct {
	ID    uint64
	Roles []string
}

func (jp *Provider) GenerateAccessToken(id uint64, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   id,
		"roles": roles,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(jp.accessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jp.secretKey))
}

func (jp *Provider) GenerateRefreshToken(id uint64, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   id,
		"roles": roles,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(jp.refreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jp.secretKey))
}

func (jp *Provider) VerifyAndParseClaims(token string) (Claims, error) {
	jwtClaims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jp.secretKey), nil
	})
	if err != nil {
		return Claims{}, err
	}

	claims := Claims{}
	claims.ID = uint64(jwtClaims["sub"].(float64))
	claims.Roles = jwtClaims["roles"].([]string)

	return claims, nil
}
