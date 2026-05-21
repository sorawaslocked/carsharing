package service

import (
	"context"

	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"
)

func userIDFromCtx(ctx context.Context) (string, error) {
	md := utils.MetadataFromCtx(ctx)
	if md.UserID == nil {
		return "", model.ErrUnauthenticated
	}
	return *md.UserID, nil
}
