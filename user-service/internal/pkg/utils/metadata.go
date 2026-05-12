package utils

import (
	"context"

	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

const (
	ctxRequestIDKey        = "x-request-id"
	ctxRequestClientIPKey  = "x-client-ip"
	ctxRequestUserIDKey    = "x-user-id"
	ctxRequestUserRolesKey = "x-user-roles"
)

type Metadata struct {
	ClientIP  string
	RequestID string
	UserID    *string
	UserRoles []model.Role
}

func MetadataFromCtx(ctx context.Context) Metadata {
	md := Metadata{}

	if clientIP, ok := ctx.Value(ctxRequestClientIPKey).(string); ok {
		md.ClientIP = clientIP
	}
	if requestID, ok := ctx.Value(ctxRequestIDKey).(string); ok {
		md.RequestID = requestID
	}
	if userID, ok := ctx.Value(ctxRequestUserIDKey).(string); ok {
		md.UserID = &userID
	}
	if rawRoles, ok := ctx.Value(ctxRequestUserRolesKey).([]string); ok {
		for _, s := range rawRoles {
			if role, err := model.RoleFromString(s); err == nil {
				md.UserRoles = append(md.UserRoles, role)
			}
		}
	}

	return md
}
