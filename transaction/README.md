# transaction

**Что:** Unit of Work поверх pgx. Прозрачная транзакция через контекст.

**Для чего:** Бизнес-логика (usecase) оборачивает несколько операций в одну транзакцию. Адаптеры не знают, работают ли они внутри транзакции или напрямую с пулом.

## Как использовать

`internal/app/app.go` — инициализация после создания пула:
```go
transaction.Init(pgPool)
```

`internal/usecase/*/scenario.go` — оборачивание нескольких операций в транзакцию:
```go
err := transaction.Wrap(ctx, func(ctx context.Context) error {
    repo.CreateUser(ctx, user)
    repo.CreateSession(ctx, session)
    return nil
})
```

`internal/adapter/postgres/*.go` — получение исполнителя (tx или пул) из контекста:
```go
conn := transaction.TryExtractTX(ctx)
conn.Exec(ctx, query, args...)
```

Паттерн: usecase управляет границами транзакции, адаптеры — прозрачны.
