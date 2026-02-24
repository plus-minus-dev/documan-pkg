# grpcserver

**Что:** Обёртка над grpc.Server с управлением жизненным циклом.

**Для чего:** Стандартный запуск gRPC-сервера в горутине. Принимает произвольные `grpc.ServerOption` для interceptor-ов.

## Как использовать

`internal/app/app.go` — создание, регистрация сервисов, запуск и shutdown:
```go
srv := grpcserver.New(cfg.GRPCPort,
    grpc.ChainUnaryInterceptor(
        logger.Interceptor(),
        otel.GRPCInterceptor(),
    ),
)

<service>_v1.Register<Service>Server(srv.Server(), controller)

srv.Start()
defer srv.GracefulStop()
```
