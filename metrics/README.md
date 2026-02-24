# metrics

**Что:** Prometheus-метрики. Готовый HTTP middleware + универсальный `Process` для произвольных операций.

**Для чего:** Автоматический сбор метрик HTTP-запросов (count, duration по method + status). `Process` — для метрик бизнес-операций (очередь, обработка, ошибки).

## Как использовать

`internal/app/app.go` — создание метрик, подключение middleware к роутеру:
```go
m := metrics.NewHTTPServer()
router.Use(metrics.NewMiddleware(m))
```

`internal/usecase/*/scenario.go` — метрики бизнес-операций:
```go
p := metrics.NewProcess("kafka_consume")
p.Current().Inc()
p.Total("process_message", metrics.Ok)
p.Duration("process_message", startTime)
```
