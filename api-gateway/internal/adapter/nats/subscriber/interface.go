package subscriber

import (
	"context"

	"carsharing/api-gateway/internal/model"
)

type UserEventHandler interface {
	OnUserCreated(ctx context.Context, userID string) error
	OnUserUpdated(ctx context.Context, userID string, isSecurityUpdate bool) error
	OnUserDeleted(ctx context.Context, userID string) error
}

type CarStatusEventHandler interface {
	OnCarStatusUpdated(ctx context.Context, event model.CarStatusUpdatedEvent) error
}
