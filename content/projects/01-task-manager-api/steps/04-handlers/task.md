---
title: "HTTP Handlers"
difficulty: medium
order: 4
file: "handler/task_handler.go"
hints:
  - "Декодируй тело запроса через json.NewDecoder(r.Body).Decode(&req)"
  - "Возвращай 201 Created для успешного создания, 204 No Content для удаления"
  - "Передавай r.Context() в методы репозитория — это позволит отменить запрос если клиент отключился"
  - "Список всегда возвращай как массив — никогда null, проверяй tasks != nil"
---

# HTTP Handlers

**HTTP Handler** — функция с сигнатурой `func(w http.ResponseWriter, r *http.Request)`. В Go handlers работают напрямую с `http.ResponseWriter` и `*http.Request` без магии фреймворков.

Ключевые принципы:
- Декодируй тело запроса, валидируй данные, вызывай бизнес-логику, пиши ответ
- Коды ответов: 200 OK, 201 Created, 204 No Content, 400 Bad Request, 404 Not Found, 500 Internal Server Error
- Всегда возвращай JSON с `Content-Type: application/json`

## Что нужно сделать

Создай `TaskHandler` в пакете `handler`.

### Интерфейс репозитория

Определи интерфейс, который использует handler. Это позволяет подменять реализацию в тестах:

```go
type TaskRepository interface {
    Create(ctx context.Context, task model.Task) (model.Task, error)
    GetByID(ctx context.Context, id int64) (model.Task, error)
    List(ctx context.Context) ([]model.Task, error)
    Update(ctx context.Context, task model.Task) (model.Task, error)
    Delete(ctx context.Context, id int64) error
}
```

### Структура и конструктор

```go
type TaskHandler struct {
    repo TaskRepository
}

func New(repo TaskRepository) *TaskHandler
```

### Методы

| Метод     | HTTP       | Успех    | Ошибка                    |
|-----------|------------|----------|---------------------------|
| `Create`  | POST /tasks | 201     | 400 (bad json / validate) |
| `GetByID` | GET /tasks/{id} | 200 | 400 (bad id), 404 (not found) |
| `List`    | GET /tasks | 200      | 500                       |
| `Update`  | PUT /tasks/{id} | 200 | 400, 404                  |
| `Delete`  | DELETE /tasks/{id} | 204 | 400, 404               |

### Вспомогательные функции

```go
func writeJSON(w http.ResponseWriter, status int, data any)
func writeError(w http.ResponseWriter, status int, message string)
```

`writeError` должна писать `{"error": "message"}`.

## Требования

- `Create` вызывает `task.Validate()` и возвращает 400 при ошибке
- ID из URL парсится через `strconv.ParseInt(r.PathValue("id"), 10, 64)`
- `r.Context()` передаётся в методы репозитория
- `List` возвращает `[]` (пустой массив JSON), не `null`
- `Delete` возвращает 204 без тела ответа

## Пример использования

```go
repo := repository.New(db)
h := handler.New(repo)

mux := http.NewServeMux()
mux.HandleFunc("POST /tasks", h.Create)
mux.HandleFunc("GET /tasks/{id}", h.GetByID)
```
