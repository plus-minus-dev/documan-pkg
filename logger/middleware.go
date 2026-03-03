package logger

import (
	"context"
	"fmt"
	"net/http"

	zlog "github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"github.com/plus-minus-dev/documan-pkg/router"
)

type ContextErrKey struct{}

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		ctx := context.WithValue(r.Context(), ContextErrKey{}, &err)
		log := zlog.Ctx(ctx)
		ctx = log.With().
			Str("transport", "http").
			Logger().
			WithContext(ctx)
		log = zlog.Ctx(ctx)

		ww := router.WriterWrapper(w)
		next.ServeHTTP(ww, r.WithContext(ctx))

		event := log.Info()

		if err != nil {
			event = log.Error().Err(err)
		}

		event.
			Int("code", ww.Code()).
			Str("method", fmt.Sprintf("%s %s", r.Method, router.ExtractPath(ctx))).
			Str("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()).
			Send()
	}

	return http.HandlerFunc(fn)
}
