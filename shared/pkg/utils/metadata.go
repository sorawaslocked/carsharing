package utils

import (
	"context"

	sharedmodel "carsharing/shared/model"
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
	UserRoles []sharedmodel.Role
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
		roles := make([]sharedmodel.Role, len(rawRoles))
		for i, r := range rawRoles {
			roles[i] = sharedmodel.Role(r)
		}
		md.UserRoles = roles
	}

	return md
}
