// Package grpcutil provides gRPC client utilities.
package grpcutil

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/plus-minus-dev/documan-pkg/ctxutil"
)

// PropagateMetadata returns a gRPC unary client interceptor that forwards
// x-request-id, x-tenant-id and x-user-id from ctxutil context values into
// outgoing gRPC metadata.
func PropagateMetadata() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var pairs []string
		if rid := ctxutil.RequestID(ctx); rid != "" {
			pairs = append(pairs, "x-request-id", rid)
		}
		if tid := ctxutil.TenantID(ctx); tid != "" {
			pairs = append(pairs, "x-tenant-id", tid)
		}
		if uid := ctxutil.UserID(ctx); uid != "" {
			pairs = append(pairs, "x-user-id", uid)
		}
		if len(pairs) > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
