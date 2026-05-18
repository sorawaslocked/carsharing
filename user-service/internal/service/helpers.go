package service

import (
	"strings"

	"carsharing/user-service/internal/model"
	"carsharing/user-service/internal/pkg/utils"
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
