package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pkgutils "carsharing/shared/pkg/utils"
)

func MetadataForwardingUnaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(attachOutgoingMetadata(ctx), method, req, reply, cc, opts...)
}

func MetadataForwardingStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return streamer(attachOutgoingMetadata(ctx), desc, cc, method, opts...)
}

func attachOutgoingMetadata(ctx context.Context) context.Context {
	md := pkgutils.MetadataFromCtx(ctx)
	var kv []string
	if md.RequestID != "" {
		kv = append(kv, "x-request-id", md.RequestID)
	}
	if md.ClientIP != "" {
		kv = append(kv, "x-client-ip", md.ClientIP)
	}
	if md.UserID != nil {
		kv = append(kv, "x-user-id", *md.UserID)
	}
	if len(md.UserRoles) > 0 {
		kv = append(kv, "x-user-roles", strings.Join(md.UserRoles, ","))
	}
	if len(kv) == 0 {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, kv...)
}
