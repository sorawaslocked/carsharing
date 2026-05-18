package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"carsharing/api-gateway/internal/config"
	"carsharing/api-gateway/internal/model"
	"carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	rdb          *redis.Client
	userProvider UserProvider
	cfg          config.CacheConfig

	log *slog.Logger
}

func NewUserCache(rdb *redis.Client, userProvider UserProvider, cfg config.CacheConfig, logger *slog.Logger) *UserCache {
	c := &UserCache{
		rdb:          rdb,
		userProvider: userProvider,
		cfg:          cfg,
	}

	c.log = log.WithComponent(logger, "redis.UserCache")

	return c
}

func (c *UserCache) Close() error {
	const method = "Close"
	logger := log.WithMethod(c.log, method)

	logger.Info("closing connection")
	err := c.rdb.Close()
	if err != nil {
		logger.Error("closing connection", log.Err(err))

		return ErrCloseFailed
	}

	return nil
}

func metadataKey(userID, field string) string {
	return fmt.Sprintf("user:%s:metadata:%s", userID, field)
}

func sessionKey(userID, deviceID string) string {
	return fmt.Sprintf("user:%s:session:%s", userID, deviceID)
}

func sessionIndexKey(userID string) string {
	return fmt.Sprintf("user:%s:sessions", userID)
}

func allMetadataKeys(userID string) []string {
	return []string{
		metadataKey(userID, "roles"),
		metadataKey(userID, "doc_verified"),
		metadataKey(userID, "email_verified"),
		metadataKey(userID, "suspended"),
	}
}

func (c *UserCache) GetRoles(ctx context.Context, userID string) ([]string, error) {
	const method = "GetRoles"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	roles, err := c.rdb.SMembers(ctx, metadataKey(userID, "roles")).Result()
	if err != nil {
		logger.Error("getting roles from redis", log.Err(err))

		return nil, ErrReadFailed
	}

	if len(roles) > 0 {
		return roles, nil
	}

	logger.Info("cache miss, restoring from provider")
	user, err := c.restore(ctx, userID)
	if err != nil {
		logger.Error("restoring from provider", log.Err(err))

		return nil, ErrWriteFailed
	}

	return user.Roles, nil
}

func (c *UserCache) IsDocumentVerified(ctx context.Context, userID string) (bool, error) {
	const method = "IsDocumentVerified"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	isDocumentVerified, err := c.getBool(ctx, userID, "doc_verified", func(u model.User) bool {
		return u.IsDocumentVerified
	})
	if err != nil {
		logger.Error("getting document verified from redis", log.Err(err))

		return false, ErrReadFailed
	}

	return isDocumentVerified, nil
}

func (c *UserCache) IsEmailVerified(ctx context.Context, userID string) (bool, error) {
	const method = "IsEmailVerified"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	isEmailVerified, err := c.getBool(ctx, userID, "email_verified", func(u model.User) bool {
		return u.IsEmailVerified
	})
	if err != nil {
		logger.Error("getting email verified from redis", log.Err(err))

		return false, ErrReadFailed
	}

	return isEmailVerified, nil
}

func (c *UserCache) IsSuspended(ctx context.Context, userID string) (bool, error) {
	const method = "IsSuspended"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	isSuspended, err := c.getBool(ctx, userID, "suspended", func(u model.User) bool {
		return u.IsSuspended
	})
	if err != nil {
		logger.Error("getting suspended from redis", log.Err(err))

		return false, ErrReadFailed
	}

	return isSuspended, nil
}

func (c *UserCache) IsSignedIn(ctx context.Context, userID, deviceID string) (bool, error) {
	const method = "IsSignedIn"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	isLoggedIn, err := c.rdb.Get(ctx, sessionKey(userID, deviceID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		logger.Error("getting session from redis", log.Err(err))

		return false, ErrReadFailed
	}

	return isLoggedIn == "1", nil
}

func (c *UserCache) SetSignedIn(ctx context.Context, userID, deviceID string, loggedIn bool) error {
	const method = "SetSignedIn"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	key := sessionKey(userID, deviceID)
	indexKey := sessionIndexKey(userID)

	if !loggedIn {
		pipe := c.rdb.Pipeline()
		pipe.Del(ctx, key)
		pipe.SRem(ctx, indexKey, deviceID)
		_, err := pipe.Exec(ctx)
		if err != nil {
			logger.Error("deleting session from redis", log.Err(err))

			return ErrDeleteFailed
		}

		return nil
	}

	pipe := c.rdb.Pipeline()
	pipe.Set(ctx, key, "1", c.cfg.SessionTTL)
	pipe.SAdd(ctx, indexKey, deviceID)
	pipe.Expire(ctx, indexKey, c.cfg.SessionTTL)
	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error("setting session in redis", log.Err(err))

		return ErrWriteFailed
	}

	return nil
}

func (c *UserCache) OnUserCreated(ctx context.Context, userID string) error {
	const method = "OnUserCreated"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := c.restore(ctx, userID)
	if err != nil {
		logger.Error("restoring from provider", log.Err(err))

		return ErrWriteFailed
	}

	return nil
}

func (c *UserCache) OnUserUpdated(ctx context.Context, userID string, isSecurityUpdate bool) error {
	const method = "OnUserUpdated"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	_, err := c.restore(ctx, userID)
	if err != nil {
		logger.Error("restoring from provider", log.Err(err))

		return ErrWriteFailed
	}

	if isSecurityUpdate {
		err = c.deleteAllSessions(ctx, userID)
		if err != nil {
			logger.Error("deleting sessions from redis", log.Err(err))

			return ErrDeleteFailed
		}
	}

	return nil
}

func (c *UserCache) OnUserDeleted(ctx context.Context, userID string) error {
	const method = "OnUserDeleted"
	logger := log.WithMethod(c.log, method)
	logger = log.WithMetadata(logger, utils.MetadataFromCtx(ctx))

	err := c.deleteMetadata(ctx, userID)
	if err != nil {
		logger.Error("deleting metadata from redis", log.Err(err))

		return ErrDeleteFailed
	}

	err = c.deleteAllSessions(ctx, userID)
	if err != nil {
		logger.Error("deleting sessions from redis", log.Err(err))

		return ErrDeleteFailed
	}

	return nil
}

func (c *UserCache) getBool(
	ctx context.Context,
	userID, field string,
	pick func(model.User) bool,
) (bool, error) {
	val, err := c.rdb.Get(ctx, metadataKey(userID, field)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	if err == nil {
		return val == "1", nil
	}

	user, err := c.restore(ctx, userID)
	if err != nil {
		return false, err
	}

	return pick(user), nil
}

func (c *UserCache) restore(ctx context.Context, userID string) (model.User, error) {
	ctx = context.WithValue(ctx, "x-user-id", userID)
	ctx = context.WithValue(ctx, "x-user-roles", []string{"user"})

	user, err := c.userProvider.Get(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	err = c.store(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (c *UserCache) store(ctx context.Context, user model.User) error {
	pipe := c.rdb.Pipeline()

	rolesKey := metadataKey(user.ID, "roles")
	pipe.Del(ctx, rolesKey)
	if len(user.Roles) > 0 {
		members := make([]interface{}, len(user.Roles))
		for i, r := range user.Roles {
			members[i] = r
		}
		pipe.SAdd(ctx, rolesKey, members...)
	}
	pipe.Expire(ctx, rolesKey, c.cfg.MetadataTTL)

	pipe.Set(ctx, metadataKey(user.ID, "doc_verified"), boolVal(user.IsDocumentVerified), c.cfg.MetadataTTL)
	pipe.Set(ctx, metadataKey(user.ID, "email_verified"), boolVal(user.IsEmailVerified), c.cfg.MetadataTTL)
	pipe.Set(ctx, metadataKey(user.ID, "suspended"), boolVal(user.IsSuspended), c.cfg.MetadataTTL)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *UserCache) deleteMetadata(ctx context.Context, userID string) error {
	if err := c.rdb.Del(ctx, allMetadataKeys(userID)...).Err(); err != nil {
		return err
	}

	return nil
}

func (c *UserCache) deleteAllSessions(ctx context.Context, userID string) error {
	indexKey := sessionIndexKey(userID)

	deviceIDs, err := c.rdb.SMembers(ctx, indexKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if len(deviceIDs) == 0 {
		return nil
	}

	keys := make([]string, 0, len(deviceIDs)+1)
	for _, deviceID := range deviceIDs {
		keys = append(keys, sessionKey(userID, deviceID))
	}
	keys = append(keys, indexKey)

	err = c.rdb.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}

	return nil
}

func boolVal(b bool) string {
	if b {
		return "1"
	}

	return "0"
}
