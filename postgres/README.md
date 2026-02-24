# postgres

**Что:** Пул соединений pgxpool + миграции goose с embedded FS.

**Для чего:** Единая точка подключения к PostgreSQL. Каждый сервис создаёт пул при старте и прогоняет миграции.

## Как использовать

`config/config.go` — добавить `postgres.Config` в структуру конфига:
```go
type Config struct {
    Postgres postgres.Config
    // ...
}
```

`internal/app/app.go` — создать пул и запустить миграции:
```go
pgPool, err := postgres.New(ctx, cfg.Postgres)
defer pgPool.Close()

pgPool.Migrate(migrationsFS, "migrations")
```

Конфигурация: `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB_NAME`. Поддерживает Docker secrets (`_FILE` суффикс — на стороне сервиса).
