# router

**Что:** HTTP-роутер chi с recovery, health-эндпоинтами и утилитами для middleware.

**Для чего:** Базовый роутер с `/live` и `/ready` (204). Утилиты нужны пакетам `metrics`, `logger`, `otel` для извлечения route pattern и статуса ответа.

## Как использовать

`internal/app/app.go` — создание роутера, подключение middleware, регистрация маршрутов:
```go
r := router.New()
r.Route("/api/v1", func(r chi.Router) {
    r.Use(otel.HTTPMiddleware, logger.Middleware, metrics.NewMiddleware(m))
    r.Post("/users", handler.CreateUser)
})

srv := httpserver.New(r, cfg.HTTP)
```

## Утилиты

- `ExtractPath(ctx)` — route pattern для метрик (`/api/v1/users`, не `/api/v1/users/123`)
- `WriterWrapper` — перехват HTTP-статуса для middleware
