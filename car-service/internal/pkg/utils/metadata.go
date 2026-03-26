package utils

import (
	"context"
	"strconv"
	"strings"
)

const (
	ctxRequestIDKey        = "x-request-id"
	ctxRequestClientIPKey  = "x-client-ip"
	ctxRequestUserIDKey    = "x-user-id"
	ctxRequestUserRoles    = "x-user-roles"
	ctxRequestUserVerified = "x-user-verified"
)

type Metadata struct {
	ClientIP     string
	RequestID    string
	UserID       string
	UserRoles    []string
	UserVerified bool
}

func MetadataFromCtx(ctx context.Context) (Metadata, bool) {
	clientIP, ok := ctx.Value(ctxRequestClientIPKey).(string)
	if !ok {
		return Metadata{}, false
	}

	requestID, ok := ctx.Value(ctxRequestIDKey).(string)
	if !ok {
		return Metadata{}, false
	}

	userID, ok := ctx.Value(ctxRequestUserIDKey).(string)
	if !ok {
		return Metadata{}, false
	}

	userRolesStr, ok := ctx.Value(ctxRequestUserRoles).(string)
	if !ok || userRolesStr == "" {
		return Metadata{}, false
	}
	userRoles := strings.Split(userRolesStr, ",")

	userVerifiedStr, ok := ctx.Value(ctxRequestUserVerified).(string)
	if !ok || userVerifiedStr == "" {
		return Metadata{}, false
	}
	userVerified, err := strconv.ParseBool(userVerifiedStr)
	if err != nil {
		return Metadata{}, false
	}

	return Metadata{
		ClientIP:     clientIP,
		RequestID:    requestID,
		UserID:       userID,
		UserRoles:    userRoles,
		UserVerified: userVerified,
	}, true
}
