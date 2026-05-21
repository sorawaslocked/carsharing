package redis

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"time"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/pkg/security"

	"github.com/redis/go-redis/v9"
)

const (
	activationCodeKeyPrefix  = "user:code:activation"
	activationCooldownPrefix = "user:code:activation:cooldown"
	codeExpiration           = 10 * time.Minute
	resendCooldown           = time.Minute
	activationCodeSymbols    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	activationCodeLength     = 6
)

type ActivationCodeCache struct {
	log *slog.Logger
	rdb *redis.Client
}

func NewActivationCodeCache(log *slog.Logger, client *redis.Client) *ActivationCodeCache {
	return &ActivationCodeCache{
		log: pkglog.WithComponent(log, "adapter.redis.ActivationCodeCache"),
		rdb: client,
	}
}

func (rc *ActivationCodeCache) key(userID string) string {
	return fmt.Sprintf("%s:%s", activationCodeKeyPrefix, userID)
}

func (rc *ActivationCodeCache) cooldownKey(userID string) string {
	return fmt.Sprintf("%s:%s", activationCooldownPrefix, userID)
}

func (rc *ActivationCodeCache) Save(ctx context.Context, userID string) (string, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(rc.log, "Save"), utils.MetadataFromCtx(ctx))

	code := createCode()

	codeHash, err := security.HashString(code)
	if err != nil {
		return "", err
	}

	if err := rc.rdb.Set(ctx, rc.key(userID), codeHash, codeExpiration).Err(); err != nil {
		log.Error("setting activation code", pkglog.Err(err))
		return "", model.ErrRedis
	}

	if err := rc.rdb.Set(ctx, rc.cooldownKey(userID), 1, resendCooldown).Err(); err != nil {
		log.Error("setting activation code cooldown", pkglog.Err(err))
	}

	return code, nil
}

func (rc *ActivationCodeCache) Get(ctx context.Context, userID string) ([]byte, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(rc.log, "Get"), utils.MetadataFromCtx(ctx))

	codeHash, err := rc.rdb.Get(ctx, rc.key(userID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, model.ErrNotFound
		}
		log.Error("getting activation code", pkglog.Err(err))
		return nil, model.ErrRedis
	}

	return codeHash, nil
}

// ResendAllowedIn returns the duration the caller must wait before requesting a
// new activation code. Returns 0 if a resend is immediately allowed.
func (rc *ActivationCodeCache) ResendAllowedIn(ctx context.Context, userID string) (time.Duration, error) {
	log := pkglog.WithMetadata(pkglog.WithMethod(rc.log, "ResendAllowedIn"), utils.MetadataFromCtx(ctx))

	ttl, err := rc.rdb.TTL(ctx, rc.cooldownKey(userID)).Result()
	if err != nil {
		log.Error("getting activation code cooldown ttl", pkglog.Err(err))
		return 0, model.ErrRedis
	}

	if ttl <= 0 {
		return 0, nil
	}

	return ttl, nil
}

func createCode() string {
	symbolRunes := []rune(activationCodeSymbols)

	var bb bytes.Buffer
	bb.Grow(activationCodeLength)
	l := uint32(len(symbolRunes))

	for i := 0; i < activationCodeLength; i++ {
		bb.WriteRune(symbolRunes[binary.BigEndian.Uint32(security.Bytes(4))%l])
	}

	return bb.String()
}
