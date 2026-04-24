---
title: "Router и запуск сервера"
difficulty: easy
order: 5
file: "server/router.go"
hints:
  - "Go 1.22+ поддерживает method routing в ServeMux: \"POST /tasks\""
  - "Параметры пути достань через r.PathValue(\"id\")"
  - "Выдели NewRouter() отдельно от Run() — так проще тестировать"
  - "Добавь GET /health endpoint — стандартная практика для любого API"
---

# Router и запуск сервера

С Go 1.22 стандартный `http.ServeMux` получил полноценный метод-роутинг. Теперь не нужны сторонние роутеры для базовых REST API.

Синтаксис Go 1.22+ ServeMux:
```go
mux.HandleFunc("POST /tasks", h.Create)        // только POST
mux.HandleFunc("GET /tasks/{id}", h.GetByID)   // параметр пути
```

## Что нужно сделать

Создай файл `server/router.go` в пакете `server`.

### Функция `NewRouter`

```go
func NewRouter(h *handler.TaskHandler) http.Handler
```

Регистрирует маршруты:

| Метод  | Путь           | Handler        |
|--------|----------------|----------------|
| GET    | /health        | inline — `{"status":"ok"}` |
| POST   | /tasks         | h.Create       |
| GET    | /tasks         | h.List         |
| GET    | /tasks/{id}    | h.GetByID      |
| PUT    | /tasks/{id}    | h.Update       |
| DELETE | /tasks/{id}    | h.Delete       |

### Функция `Run`

```go
func Run(addr string, handler http.Handler) error
```

Просто оборачивает `http.ListenAndServe(addr, handler)`.

## Требования

- `NewRouter` возвращает `http.Handler` (не `*http.ServeMux`) — это позволит позже добавить middleware
- `/health` возвращает `{"status": "ok"}` с `Content-Type: application/json`
- Используй `http.NewServeMux()` (не `http.DefaultServeMux`)
- `NewRouter` и `Run` — разные функции, это упрощает тестирование

## Пример использования

```go
repo := repository.New(db)
h := handler.New(repo)
router := server.NewRouter(h)

log.Fatal(server.Run(":8080", router))
```
