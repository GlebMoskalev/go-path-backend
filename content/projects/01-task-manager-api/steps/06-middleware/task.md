---
title: "Middleware: логирование и recovery"
difficulty: medium
order: 6
file: "middleware/middleware.go"
hints:
  - "Middleware имеет сигнатуру func(http.Handler) http.Handler"
  - "Для замера времени: start := time.Now(), потом time.Since(start)"
  - "Recovery ловит panic через defer + recover(), возвращает 500"
  - "Чтобы залогировать статус ответа — оберни ResponseWriter в свою структуру"
---

# Middleware: логирование и recovery

**Middleware** — функция, которая оборачивает http.Handler и добавляет поведение до/после обработки запроса. Это стандартный паттерн для сквозной функциональности: логирование, аутентификация, rate limiting, recovery от паник.

Сигнатура middleware в Go:
```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // до
        next.ServeHTTP(w, r)
        // после
    })
}
```

## Что нужно сделать

Создай три компонента в пакете `middleware`.

### 1. `responseWriter` — обёртка для перехвата статуса

```go
type responseWriter struct {
    http.ResponseWriter
    status int
}
```

- Встраивает `http.ResponseWriter`
- Переопределяет `WriteHeader(code int)` — сохраняет `code` в поле `status`
- По умолчанию `status = http.StatusOK`

### 2. `Logger` — логирование запросов

```go
func Logger(next http.Handler) http.Handler
```

Логирует каждый запрос в формате:
```
METHOD /path → status elapsed
```
Например: `GET /tasks → 200 1.2ms`

Используй `log.Printf` из стандартной библиотеки.

### 3. `Recovery` — защита от паник

```go
func Recovery(next http.Handler) http.Handler
```

Ловит `panic` через `defer + recover()`. При панике:
- Возвращает 500 Internal Server Error
- Логирует панику через `log.Printf`

### 4. `Chain` — цепочка middleware

```go
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler
```

Применяет middleware в порядке передачи: первый в списке — самый внешний.

## Требования

- `Recovery` должен работать даже если handler уже вызвал `WriteHeader`
- `Logger` использует `responseWriter` чтобы залогировать реальный статус ответа
- `Chain` применяет middleware справа налево (чтобы первый в списке был внешним)

## Пример использования

```go
router := server.NewRouter(h)
handler := middleware.Chain(router,
    middleware.Logger,
    middleware.Recovery,
)
```
