# redis

**Что:** Клиент go-redis с graceful shutdown.

**Для чего:** Подключение к Redis для кэша, очередей, сессий.

## Как использовать

`config/config.go` — добавить в структуру конфига:
```go
type Config struct {
    Redis redis.Config
    // ...
}
```

`internal/app/app.go` — создание клиента и shutdown:
```go
client := redis.New(cfg.Redis)
defer client.Close()
```

`internal/adapter/redis/*.go` — использование клиента в адаптерах:
```go
client.Set(ctx, "key", "value", ttl)
val, err := client.Get(ctx, "key").Result()
```

Конфигурация: `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`.
