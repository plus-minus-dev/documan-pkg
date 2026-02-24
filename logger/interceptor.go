package logger

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var err error
	ctx = context.WithValue(ctx, ContextErrKey{}, &err)

	event := log.Info()

	resp, err := handler(ctx, req)
	if err != nil {
		event = log.Error().Err(err)
	}

	event.
		Str("code", status.Code(err).String()).
		Str("grpc_method", info.FullMethod).
		Send()

	return resp, err
}
