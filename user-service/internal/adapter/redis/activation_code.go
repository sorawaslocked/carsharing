package redis

import (
	"bytes"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"

	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/pkg/security"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	activationCodeKeyPrefix = "user:code:activation"
	codeExpiration          = 10 * time.Minute
	activationCodeSymbols   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	activationCodeLength    = 6
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
