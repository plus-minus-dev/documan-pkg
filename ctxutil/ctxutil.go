// Package ctxutil provides typed context keys for cross-cutting concerns
// (request ID, tenant ID) shared across HTTP and gRPC layers.
package ctxutil

import "context"

type requestIDKey struct{}
type tenantIDKey struct{}
type userIDKey struct{}

// WithRequestID stores request ID in context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID extracts request ID from context (empty string if absent).
func RequestID(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey{}).(string)
	return v
}

// WithTenantID stores tenant ID in context.
func WithTenantID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, tenantIDKey{}, id)
}

// TenantID extracts tenant ID from context (empty string if absent).
func TenantID(ctx context.Context) string {
	v, _ := ctx.Value(tenantIDKey{}).(string)
	return v
}

// WithUserID stores user ID in context.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

// UserID extracts user ID from context (empty string if absent).
func UserID(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey{}).(string)
	return v
}
