package redis

import (
	"context"

	"carsharing/api-gateway/internal/model"
)

type UserProvider interface {
	Get(ctx context.Context, id string) (model.User, error)
}
