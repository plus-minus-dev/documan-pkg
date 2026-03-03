# documan-pkg — контекст для AI-агента

> Код — источник истины.

## Роль

Общие инфраструктурные пакеты для сервисов DocuMan.
Не деплоится самостоятельно — подключается как Go-модуль (`go get`).

## Структура

```
postgres/       — pgxpool + goose миграции
transaction/    — Unit of Work через контекст
logger/         — zerolog + HTTP/gRPC middleware
grpcserver/     — обёртка grpc.Server с lifecycle
httpserver/     — обёртка http.Server с graceful shutdown
router/         — chi + health-эндпоинты
render/         — JSON-ответы для HTTP-контроллеров
otel/           — OpenTelemetry: OTLP, middleware, спаны
metrics/        — Prometheus: HTTP middleware + Process
redis/          — go-redis с graceful shutdown
ctxutil/        — утилиты контекста
grpcutil/       — утилиты gRPC
contract/       — контрактные хелперы
services/       — proto-контракты и сгенерированный код (см. ниже)
```

## Контрактный поток (proto → generate → publish)

Proto-контракты всех gRPC-сервисов живут здесь, в `services/proto/`.

### Структура

```
services/
  buf.yaml          — lint (STANDARD) + breaking (FILE)
  buf.gen.yaml      — protoc-gen-go + protoc-gen-go-grpc
  proto/
    auth/v1/        — AuthService (12 rpc)
    core/v1/        — CoreService (commands, queries, etl, push)
    connector_ms/v1/— ConnectorMsService
    common/v1/      — общие типы (Pagination и т.д.)
  gen/
    auth/v1/        — сгенерированный Go-код
    core/v1/
    connector_ms/v1/
    common/v1/
```

### Порядок обновления контракта

1. Правь proto в `services/proto/<service>/v1/`
2. `cd services && buf lint` — проверка стиля
3. `cd services && buf breaking --against '.git#subdir=services'` — проверка обратной совместимости
4. `cd services && buf generate` — генерация Go-кода в `services/gen/`
5. `git add . && git commit && git push` (или tag)
6. Потребители обновляют: `go get github.com/plus-minus-dev/documan-pkg@<commit/tag>`

### Потребители

| Сервис | Импортирует |
|--------|------------|
| documan-auth | `services/gen/auth/v1` |
| documan-core | `services/gen/core/v1`, `services/gen/common/v1` |
| documan-bff | `services/gen/auth/v1`, `services/gen/core/v1`, `services/gen/connector_ms/v1` |
| connector-ms | `services/gen/connector_ms/v1`, `services/gen/core/v1` |

### Правила

- **Никаких `replace` или `go.work`** в сервисах для documan-pkg — только version bump через `go get`.
- **Breaking changes**: buf breaking проверяет по FILE policy. Если breaking неизбежен — координируй обновление всех потребителей.
- **Генерация**: только через `buf generate`, не вручную protoc.
