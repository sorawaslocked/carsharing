package subscriber

import (
	"context"
)

type UserEventHandler interface {
	OnUserCreated(ctx context.Context, userID string) error
	OnUserUpdated(ctx context.Context, userID string, isSecurityUpdate bool) error
	OnUserDeleted(ctx context.Context, userID string) error
}
