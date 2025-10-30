package interceptor

type JwtProvider interface {
	GenerateAccessToken(id uint64, roles []string) (string, error)
	GenerateRefreshToken(id uint64, roles []string) (string, error)
	VerifyAndParseClaims(token string) (uint64, []string, error)
}
