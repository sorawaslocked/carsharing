package service

import (
	"strings"

	"context"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/utils"
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
