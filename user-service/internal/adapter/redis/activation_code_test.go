package redis_test

import (
	"context"
	"io"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	rcache "carsharing/user-service/internal/adapter/redis"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/pkg/security"
)

// Key templates mirroring the private constants in activation_code.go.
const (
	codeKeyPrefix     = "user:code:activation:"
	cooldownKeyPrefix = "user:code:activation:cooldown:"
	codeExpiration    = 10 * time.Minute
	resendCooldown    = time.Minute

	testUserID = "11111111-1111-1111-1111-111111111111"
)

var codePattern = regexp.MustCompile(`^[0-9A-Z]{6}$`)

func newCache(t *testing.T) (*rcache.ActivationCodeCache, *miniredis.Miniredis) {
	t.Helper()
	s := miniredis.RunT(t)
	client := goredis.NewClient(&goredis.Options{Addr: s.Addr()})
	t.Cleanup(func() { client.Close() })
	return rcache.NewActivationCodeCache(slog.New(slog.NewTextHandler(io.Discard, nil)), client), s
}

// --- Save ---

func TestActivationCodeCache_Save_ReturnsAlphanumericCode(t *testing.T) {
	cache, _ := newCache(t)

	code, err := cache.Save(context.Background(), testUserID)

	require.NoError(t, err)
	assert.Regexp(t, codePattern, code)
}

func TestActivationCodeCache_Save_SetsCodeKeyWithTTL(t *testing.T) {
	cache, s := newCache(t)

	_, err := cache.Save(context.Background(), testUserID)
	require.NoError(t, err)

	ttl := s.TTL(codeKeyPrefix + testUserID)
	assert.Greater(t, ttl, 9*time.Minute)
	assert.LessOrEqual(t, ttl, codeExpiration)
}

func TestActivationCodeCache_Save_SetsCooldownKeyWithTTL(t *testing.T) {
	cache, s := newCache(t)

	_, err := cache.Save(context.Background(), testUserID)
	require.NoError(t, err)

	ttl := s.TTL(cooldownKeyPrefix + testUserID)
	assert.Greater(t, ttl, 30*time.Second)
	assert.LessOrEqual(t, ttl, resendCooldown)
}

func TestActivationCodeCache_Save_OverwritesPreviousCode(t *testing.T) {
	cache, _ := newCache(t)
	ctx := context.Background()

	first, err := cache.Save(ctx, testUserID)
	require.NoError(t, err)

	second, err := cache.Save(ctx, testUserID)
	require.NoError(t, err)

	// The two codes should differ (with overwhelming probability).
	assert.NotEqual(t, first, second)

	// Get should return a hash matching the second code, not the first.
	hash, err := cache.Get(ctx, testUserID)
	require.NoError(t, err)
	assert.NoError(t, security.CheckStringHash(second, hash))
	assert.Error(t, security.CheckStringHash(first, hash))
}

// --- Get ---

func TestActivationCodeCache_Get_ReturnsHashMatchingCode(t *testing.T) {
	cache, _ := newCache(t)
	ctx := context.Background()

	code, err := cache.Save(ctx, testUserID)
	require.NoError(t, err)

	hash, err := cache.Get(ctx, testUserID)

	require.NoError(t, err)
	assert.NoError(t, security.CheckStringHash(code, hash))
}

func TestActivationCodeCache_Get_NotFound_WhenKeyAbsent(t *testing.T) {
	cache, _ := newCache(t)

	_, err := cache.Get(context.Background(), testUserID)

	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestActivationCodeCache_Get_NotFound_AfterCodeExpiry(t *testing.T) {
	cache, s := newCache(t)
	ctx := context.Background()

	_, err := cache.Save(ctx, testUserID)
	require.NoError(t, err)

	s.FastForward(codeExpiration + time.Second)

	_, err = cache.Get(ctx, testUserID)
	assert.ErrorIs(t, err, model.ErrNotFound)
}

// --- ResendAllowedIn ---

func TestActivationCodeCache_ResendAllowedIn_ZeroWhenNoCooldown(t *testing.T) {
	cache, _ := newCache(t)

	wait, err := cache.ResendAllowedIn(context.Background(), testUserID)

	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), wait)
}

func TestActivationCodeCache_ResendAllowedIn_PositiveAfterSave(t *testing.T) {
	cache, _ := newCache(t)
	ctx := context.Background()

	_, err := cache.Save(ctx, testUserID)
	require.NoError(t, err)

	wait, err := cache.ResendAllowedIn(ctx, testUserID)

	require.NoError(t, err)
	assert.Greater(t, wait, time.Duration(0))
	assert.LessOrEqual(t, wait, resendCooldown)
}

func TestActivationCodeCache_ResendAllowedIn_ZeroAfterCooldownExpiry(t *testing.T) {
	cache, s := newCache(t)
	ctx := context.Background()

	_, err := cache.Save(ctx, testUserID)
	require.NoError(t, err)

	s.FastForward(resendCooldown + time.Second)

	wait, err := cache.ResendAllowedIn(ctx, testUserID)

	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), wait)
}

func TestActivationCodeCache_CodeRemainsValid_AfterCooldownExpiry(t *testing.T) {
	cache, s := newCache(t)
	ctx := context.Background()

	code, err := cache.Save(ctx, testUserID)
	require.NoError(t, err)

	// Cooldown expires but the code (10 min TTL) should still be there.
	s.FastForward(resendCooldown + time.Second)

	hash, err := cache.Get(ctx, testUserID)
	require.NoError(t, err)
	assert.NoError(t, security.CheckStringHash(code, hash))
}

func TestActivationCodeCache_CooldownAndCode_IndependentTTLs(t *testing.T) {
	cache, s := newCache(t)
	ctx := context.Background()

	_, err := cache.Save(ctx, testUserID)
	require.NoError(t, err)

	codeTTL := s.TTL(codeKeyPrefix + testUserID)
	cooldownTTL := s.TTL(cooldownKeyPrefix + testUserID)

	assert.Greater(t, codeTTL, cooldownTTL, "code TTL should be longer than cooldown TTL")
}
