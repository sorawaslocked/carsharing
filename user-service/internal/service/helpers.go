package service

import (
	"strings"

	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"
	"context"
)

func uncapitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func userIDFromCtx(ctx context.Context) (string, error) {
	md := utils.MetadataFromCtx(ctx)
	if md.UserID == nil {
		return "", model.ErrUnauthenticated
	}
	return *md.UserID, nil
}
