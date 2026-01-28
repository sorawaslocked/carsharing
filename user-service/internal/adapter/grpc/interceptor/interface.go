package interceptor

import "time"

type JwtProvider interface {
	GenerateAccessToken(id uint64, roles []string) (string, time.Time, error)
	GenerateRefreshToken(id uint64, roles []string) (string, time.Time, error)
	VerifyAndParseClaims(token string) (uint64, []string, error)
}
