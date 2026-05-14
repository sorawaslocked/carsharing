package utils

import (
	"context"
	"strings"

	"github.com/sorawaslocked/car-rental-car-service/internal/model"
)

type contextKey string

const (
	ctxKeyRequestID contextKey = "x-request-id"
	ctxKeyClientIP  contextKey = "x-client-ip"
	ctxKeyUserID    contextKey = "x-user-id"
	ctxKeyUserRoles contextKey = "x-user-roles"
)

type Metadata struct {
	ClientIP  string
	RequestID string
	UserID    *string
	UserRoles []model.Role
}

func SetMetadata(ctx context.Context, requestID, clientIP, userID, userRoles string) context.Context {
	ctx = context.WithValue(ctx, ctxKeyRequestID, requestID)
	ctx = context.WithValue(ctx, ctxKeyClientIP, clientIP)
	ctx = context.WithValue(ctx, ctxKeyUserID, userID)
	ctx = context.WithValue(ctx, ctxKeyUserRoles, userRoles)
	return ctx
}

func MetadataFromCtx(ctx context.Context) Metadata {
	md := Metadata{}

	if v, ok := ctx.Value(ctxKeyClientIP).(string); ok {
		md.ClientIP = v
	}
	if v, ok := ctx.Value(ctxKeyRequestID).(string); ok {
		md.RequestID = v
	}
	if v, ok := ctx.Value(ctxKeyUserID).(string); ok && v != "" {
		md.UserID = &v
	}
	if v, ok := ctx.Value(ctxKeyUserRoles).(string); ok && v != "" {
		for _, s := range strings.Split(v, ",") {
			if s = strings.TrimSpace(s); s != "" {
				md.UserRoles = append(md.UserRoles, model.Role(s))
			}
		}
	}

	return md
}
