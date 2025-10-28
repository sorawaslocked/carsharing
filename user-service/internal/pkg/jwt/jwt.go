package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Config struct {
	SecretKey       string        `env:"JWT_SECRET_KEY" env-required:"true"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"JWT_ACCESS_TOKEN_TTL" env-default:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"JWT_REFRESH_TOKEN_TTL" env-default:"24h"`
}

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

func (jp *Provider) VerifyAndParseClaims(token string) (uint64, []string, error) {
	jwtClaims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jp.secretKey), nil
	})
	if err != nil {
		return 0, nil, err
	}

	id := uint64(jwtClaims["sub"].(float64))
	roles := jwtClaims["roles"].([]interface{})

	roleStrings := make([]string, len(roles))
	for i, v := range roles {
		roleStrings[i] = v.(string)
	}

	return id, roleStrings, nil
}
