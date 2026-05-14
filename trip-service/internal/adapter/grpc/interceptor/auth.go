package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sorawaslocked/car-rental-trip-service/internal/pkg/utils"
)

type AccessPolicy int

const (
	PolicyPublic AccessPolicy = iota
	PolicyAuthenticated
)

var methodPolicies = map[string]AccessPolicy{
	"/service.trip.HealthService/Health": PolicyPublic,
}

func AuthUnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctx = extractMetadata(ctx)
	if err := enforce(ctx, policyFor(info.FullMethod)); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func AuthStreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := extractMetadata(ss.Context())
	if err := enforce(ctx, policyFor(info.FullMethod)); err != nil {
		return err
	}
	return handler(srv, &wrappedStream{ss, ctx})
}

func policyFor(method string) AccessPolicy {
	if p, ok := methodPolicies[method]; ok {
		return p
	}
	return PolicyAuthenticated
}

func enforce(ctx context.Context, policy AccessPolicy) error {
	if policy == PolicyPublic {
		return nil
	}
	if utils.MetadataFromCtx(ctx).UserID == nil {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return nil
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
