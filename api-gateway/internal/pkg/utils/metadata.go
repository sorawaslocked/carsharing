package utils

import (
	"context"
)

const (
	ctxRequestIDKey       = "x-request-id"
	ctxRequestClientIPKey = "x-client-ip"
	ctxRequestUserIDKey   = "x-user-id"
)

type Metadata struct {
	ClientIP  string
	RequestID string
	UserID    *string
}

func MetadataFromCtx(ctx context.Context) Metadata {
	md := Metadata{}

	clientIP, ok := ctx.Value(ctxRequestClientIPKey).(string)
	if ok {
		md.ClientIP = clientIP
	}
	requestID, ok := ctx.Value(ctxRequestIDKey).(string)
	if ok {
		md.RequestID = requestID
	}
	userID, ok := ctx.Value(ctxRequestUserIDKey).(string)
	if ok {
		md.UserID = &userID
	}

	return md
}
