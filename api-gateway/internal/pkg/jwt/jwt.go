package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(
	secret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Manager {
	return &Manager{
		secret:          secret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (m *Manager) GenerateAccessToken(userID string) (token string, exp time.Time, err error) {
	now := time.Now()
	exp = now.Add(m.accessTokenTTL)

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = tokenWithClaims.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, exp, nil
}

func (m *Manager) GenerateRefreshToken(userID string) (token string, exp time.Time, err error) {
	now := time.Now()
	exp = now.Add(m.refreshTokenTTL)

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = tokenWithClaims.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, exp, nil
}

func (m *Manager) ParseToken(token string) (string, error) {
	var tokenWithClaims *jwt.Token
	tokenWithClaims, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return m.secret, nil
	})
	if err != nil {
		return "", ErrInvalidToken
	}

	claims, ok := tokenWithClaims.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}
