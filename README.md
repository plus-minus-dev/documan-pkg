# documan-pkg

Общие инфраструктурные пакеты для сервисов DocuMan.

| Пакет | Назначение |
|-------|-----------|
| [postgres](postgres/) | Пул соединений pgxpool + миграции goose |
| [transaction](transaction/) | Unit of Work — прозрачная транзакция через контекст |
| [logger](logger/) | Структурированный логгер zerolog + HTTP/gRPC middleware |
| [grpcserver](grpcserver/) | Обёртка над grpc.Server с lifecycle |
| [httpserver](httpserver/) | Обёртка над http.Server с таймаутами и graceful shutdown |
| [router](router/) | HTTP-роутер chi + health-эндпоинты + утилиты для middleware |
| [render](render/) | JSON-ответы и обработка ошибок для HTTP-контроллеров |
| [otel](otel/) | OpenTelemetry трейсинг: OTLP-экспорт, HTTP/gRPC middleware, ручные спаны |
| [metrics](metrics/) | Prometheus-метрики: HTTP middleware + универсальный Process |
| [redis](redis/) | Клиент go-redis с graceful shutdown |
