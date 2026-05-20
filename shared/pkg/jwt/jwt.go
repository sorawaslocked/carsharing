package jwt

import (
	"carsharing/shared/pkg/utils"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	pkglog "carsharing/shared/pkg/log"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret          string        `yaml:"secret" env:"JWT_SECRET" env-required:"true"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"JWT_ACCESS_TOKEN_TTL" env-default:"30m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"JWT_REFRESH_TOKEN_TTL" env-default:"720h"`
}

type Manager struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	log *slog.Logger
}

func NewManager(cfg Config, log *slog.Logger) *Manager {
	m := &Manager{
		secret:          cfg.Secret,
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}

	m.log = pkglog.WithComponent(log, "jwt.Manager")

	return m
}

func (m *Manager) GenerateAccessToken(ctx context.Context, userID string) (token string, exp time.Time, err error) {
	log := pkglog.WithMethod(m.log, "GenerateAccessToken")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	now := time.Now()
	exp = now.Add(m.accessTokenTTL)
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = tokenWithClaims.SignedString([]byte(m.secret))
	if err != nil {
		log.Error("generating access token", pkglog.Err(err))

		return "", time.Time{}, ErrTokenGenerationFailed
	}

	return token, exp, nil
}

func (m *Manager) GenerateRefreshToken(ctx context.Context, userID string) (token string, exp time.Time, err error) {
	log := pkglog.WithMethod(m.log, "GenerateRefreshToken")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	now := time.Now()
	exp = now.Add(m.refreshTokenTTL)
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = tokenWithClaims.SignedString([]byte(m.secret))
	if err != nil {
		log.Error("generating refresh token", pkglog.Err(err))

		return "", time.Time{}, ErrTokenGenerationFailed
	}

	return token, exp, nil
}

func (m *Manager) ParseToken(ctx context.Context, token string) (userID string, exp time.Time, err error) {
	log := pkglog.WithMethod(m.log, "ParseToken")
	log = pkglog.WithMetadata(log, utils.MetadataFromCtx(ctx))

	tokenWithClaims, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", time.Time{}, ErrExpiredToken
		}

		log.Warn("parsing token", pkglog.Err(err))

		return "", time.Time{}, ErrInvalidToken
	}

	claims, ok := tokenWithClaims.Claims.(*jwt.MapClaims)
	if !ok {
		log.Error("unexpected claims type in token")

		return "", time.Time{}, ErrInvalidToken
	}

	subject, err := claims.GetSubject()
	if err != nil {
		log.Error("getting subject from token claims", pkglog.Err(err))

		return "", time.Time{}, ErrInvalidToken
	}

	expTime, err := claims.GetExpirationTime()
	if err != nil {
		log.Error("getting expiry from token claims", pkglog.Err(err))

		return "", time.Time{}, ErrInvalidToken
	}

	return subject, expTime.Time, nil
}
