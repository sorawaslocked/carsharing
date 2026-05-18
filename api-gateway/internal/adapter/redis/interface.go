package redis

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type UserProvider interface {
	Get(ctx context.Context, id string) (model.User, error)
}
