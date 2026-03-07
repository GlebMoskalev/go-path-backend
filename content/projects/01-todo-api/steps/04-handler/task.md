---
title: "HTTP обработчики"
description: "Реализуйте REST API эндпоинты для управления задачами"
order: 4
difficulty: hard
file: "handler/todo.go"
hints:
  - "Используйте http.NewServeMux() с Go 1.22+ синтаксисом: mux.HandleFunc(\"GET /todos\", h.GetAll)"
  - "Для получения path-параметра используйте r.PathValue(\"id\")"
  - "В Create возвращайте статус 201 (http.StatusCreated), в Delete — 204 (http.StatusNoContent)"
  - "Для проверки ошибки not found используйте errors.Is(err, storage.ErrNotFound)"
  - "Не забудьте вспомогательные функции writeJSON и writeError для отправки JSON-ответов"
---

# HTTP обработчики (Handler)

**Handler (обработчик)** — это слой, который принимает HTTP-запросы, вызывает сервис и формирует HTTP-ответы. Это «лицо» вашего API.

## REST API дизайн

Наш Todo API следует REST-конвенциям:

| Метод   | Путь          | Действие         | Статус ответа |
|---------|---------------|------------------|---------------|
| GET     | `/todos`      | Список задач     | 200 OK        |
| POST    | `/todos`      | Создать задачу   | 201 Created   |
| GET     | `/todos/{id}` | Получить задачу  | 200 OK        |
| PUT     | `/todos/{id}` | Обновить задачу  | 200 OK        |
| DELETE  | `/todos/{id}` | Удалить задачу   | 204 No Content|

## Что нужно реализовать

### Структура `TodoHandler`

```go
type TodoHandler struct {
    service *service.TodoService
}
```

### Конструктор

```go
func NewTodoHandler(svc *service.TodoService) *TodoHandler
```

### Метод `Routes`

Возвращает `http.Handler` с зарегистрированными маршрутами:

```go
func (h *TodoHandler) Routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /todos", h.GetAll)
    mux.HandleFunc("POST /todos", h.Create)
    mux.HandleFunc("GET /todos/{id}", h.GetByID)
    mux.HandleFunc("PUT /todos/{id}", h.Update)
    mux.HandleFunc("DELETE /todos/{id}", h.Delete)
    return mux
}
```

### Методы-обработчики

- **Create** — декодирует JSON, проверяет что Title не пуст, вызывает сервис, возвращает 201
- **GetAll** — вызывает сервис, возвращает JSON-массив
- **GetByID** — парсит ID из пути, обрабатывает not found (404)
- **Update** — парсит ID, декодирует JSON, вызывает сервис
- **Delete** — парсит ID, вызывает сервис, возвращает 204

### Вспомогательные функции

```go
func writeJSON(w http.ResponseWriter, code int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
    writeJSON(w, code, map[string]string{"error": msg})
}
```

## Go 1.22+ маршрутизация

Начиная с Go 1.22, `http.ServeMux` поддерживает методы и path-параметры:

```go
// Метод + путь
mux.HandleFunc("GET /todos/{id}", handler)

// Получить параметр из пути
id := r.PathValue("id")
```
