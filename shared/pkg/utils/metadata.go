package utils

import (
	"context"
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
	UserRoles []string
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
	if userRoles, ok := ctx.Value(ctxRequestUserRolesKey).([]string); ok {
		md.UserRoles = userRoles
	}

	return md
}
