package handler

import (
	"context"

	"github.com/sorawaslocked/car-rental-api-gateway/internal/model"
)

type UserEventHandler interface {
	OnUserCreated(ctx context.Context, userID string) error
	OnUserUpdated(ctx context.Context, userID string, isSecurityUpdate bool) error
	OnUserDeleted(ctx context.Context, userID string) error
}

type DocumentEventHandler interface {
	OnDocumentAnalyzed(ctx context.Context, event model.DocumentAnalyzedEvent) error
}

type CarStatusEventHandler interface {
	OnCarStatusUpdated(ctx context.Context, event model.CarStatusUpdated) error
}
