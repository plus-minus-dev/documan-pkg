# logger

**Что:** Структурированный логгер zerolog + HTTP middleware + gRPC interceptor.

**Для чего:** Единый формат логов во всех сервисах. Middleware автоматически логирует входящие запросы (метод, путь, статус, trace ID).

## Как использовать

`config/config.go` — добавить в структуру конфига:
```go
type Config struct {
    Logger logger.Config
    // ...
}
```

`cmd/app/main.go` — инициализация до запуска приложения:
```go
logger.Init(cfg.Logger)
```

`internal/app/app.go` — подключить middleware или interceptor при создании сервера:
```go
// HTTP
router.Use(logger.Middleware)

// gRPC
grpc.NewServer(grpc.UnaryInterceptor(logger.Interceptor()))
```

`PrettyConsole: true` — человекочитаемый вывод для локальной разработки. `false` — JSON для продакшена.
