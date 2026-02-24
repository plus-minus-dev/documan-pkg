# render

**Что:** Формирование HTTP-ответов: JSON и ошибки.

**Для чего:** Единообразные ответы из HTTP-контроллеров. `Error()` записывает ошибку в контекст — logger middleware подхватит и залогирует.

## Как использовать

`internal/controller/http/v1/*.go` — в HTTP-хендлерах для формирования ответов:
```go
// Успех
render.JSON(w, responseDTO, http.StatusOK)

// Ошибка
render.Error(ctx, err, "failed to create user")
render.JSON(w, errorResponse, http.StatusBadRequest)
```
