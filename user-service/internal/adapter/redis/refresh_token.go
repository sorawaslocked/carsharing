package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"time"
)

const (
	refreshTokenKeyPrefix = "user:token:refresh"
)

type SessionRedisCache struct {
	rdb             *redis.Client
	refreshTokenTTL time.Duration
}

func NewSessionRedisCache(client *redis.Client, refreshTokenTTL time.Duration) *SessionRedisCache {
	return &SessionRedisCache{
		rdb:             client,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (rc *SessionRedisCache) key(userID uint64) string {
	return fmt.Sprintf("%s:%d", refreshTokenKeyPrefix, userID)
}

func (rc *SessionRedisCache) Save(ctx context.Context, userID uint64) error {
	err := rc.rdb.Set(ctx, rc.key(userID), true, rc.refreshTokenTTL).Err()
	if err != nil {
		return model.ErrRedis
	}

	return nil
}

func (rc *SessionRedisCache) Exists(ctx context.Context, userID uint64) (bool, error) {
	_, err := rc.rdb.Get(ctx, rc.key(userID)).Bool()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, model.ErrNotFound
		}

		return false, model.ErrRedis
	}

	return true, nil
}

func (rc *SessionRedisCache) Delete(ctx context.Context, userID uint64) error {
	err := rc.rdb.Del(ctx, rc.key(userID)).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.ErrNotFound
		}

		return model.ErrRedis
	}

	return nil
}
