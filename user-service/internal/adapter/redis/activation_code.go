package redis

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/security"
	"time"
)

const (
	activationCodeKeyPrefix = "user:code:activation"
	codeExpiration          = 10 * time.Minute
	activationCodeSymbols   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	activationCodeLength    = 6
)

type ActivationCodeRedisCache struct {
	rdb *redis.Client
}

func NewActivationCodeRedisCache(client *redis.Client) *ActivationCodeRedisCache {
	return &ActivationCodeRedisCache{
		rdb: client,
	}
}

func (rc *ActivationCodeRedisCache) key(userID uint64) string {
	return fmt.Sprintf("%s:%d", activationCodeKeyPrefix, userID)
}

func (rc *ActivationCodeRedisCache) Save(ctx context.Context, userID uint64) (string, error) {
	code := createCode()

	codeHash, err := security.HashString(code)
	if err != nil {
		return "", err
	}

	err = rc.rdb.Set(ctx, rc.key(userID), codeHash, codeExpiration).Err()

	return code, model.ErrRedis
}

func (rc *ActivationCodeRedisCache) Get(ctx context.Context, userID uint64) ([]byte, error) {
	codeHash, err := rc.rdb.Get(ctx, rc.key(userID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, model.ErrNotFound
		}

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
