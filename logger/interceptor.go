package logger

import (
	"context"

	zlog "github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var err error
	ctx = context.WithValue(ctx, ContextErrKey{}, &err)
	log := zlog.Ctx(ctx)
	ctx = log.With().
		Str("transport", "grpc").
		Str("grpc_method", info.FullMethod).
		Logger().
		WithContext(ctx)
	log = zlog.Ctx(ctx)

	event := log.Info()

	resp, err := handler(ctx, req)
	if err != nil {
		event = log.Error().Err(err)
	}

	event.
		Str("code", status.Code(err).String()).
		Str("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()).
		Send()

	return resp, err
}
