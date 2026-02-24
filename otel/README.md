# otel

**Что:** OpenTelemetry трейсинг. OTLP-экспорт через gRPC. HTTP middleware, gRPC interceptor, ручные спаны.

**Для чего:** Распределённый трейсинг между сервисами. Если `OTEL_ENDPOINT` пуст — автоматически переходит в noop-режим (без накладных расходов).

## Как использовать

`config/config.go` — добавить в структуру конфига:
```go
type Config struct {
    Otel otel.Config
    // ...
}
```

`internal/app/app.go` — инициализация и подключение middleware/interceptor:
```go
otel.Init(ctx, cfg.Otel)
defer otel.Close()

// HTTP
router.Use(otel.HTTPMiddleware)

// gRPC
grpc.ChainUnaryInterceptor(otel.GRPCInterceptor())
```

`internal/usecase/*/scenario.go`, `internal/adapter/postgres/*.go` — ручные спаны:
```go
ctx, span := tracer.Start(ctx, "CreateUser")
defer span.End()
```

Конфигурация: `OTEL_ENDPOINT`, `OTEL_NAMESPACE`, `OTEL_INSTANCE_ID`, `OTEL_RATIO` (sampling).
