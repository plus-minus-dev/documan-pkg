# httpserver

**Что:** Обёртка над http.Server с таймаутами (20s read/write) и graceful shutdown (25s).

**Для чего:** Стандартный запуск HTTP-сервера в горутине. Принимает любой `http.Handler`.

## Как использовать

`config/config.go` — добавить в структуру конфига:
```go
type Config struct {
    HTTP httpserver.Config
    // ...
}
```

`internal/app/app.go` — создание и shutdown:
```go
srv := httpserver.New(router, cfg.HTTP)
defer srv.Close()
```
