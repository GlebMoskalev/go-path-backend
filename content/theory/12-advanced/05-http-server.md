---
title: "HTTP-сервер"
description: "net/http, Handler интерфейс, ServeMux, middleware, graceful shutdown"
order: 5
---

# HTTP-сервер

Пакет `net/http` из стандартной библиотеки содержит всё необходимое для production-ready HTTP-сервера.

## Минимальный сервер

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })

    http.ListenAndServe(":8080", nil)
}
```

---

## Интерфейс http.Handler

Всё основано на одном интерфейсе:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

`http.HandlerFunc` — адаптер для превращения функции в `Handler`:

```go
type HandlerFunc func(ResponseWriter, *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) { f(w, r) }
```

Собственный тип-обработчик:

```go
type UserHandler struct {
    db *sql.DB
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.getUser(w, r)
    case http.MethodPost:
        h.createUser(w, r)
    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
}
```

---

## http.ServeMux — маршрутизатор

```go
mux := http.NewServeMux()

// Go 1.22+: метод и путь в паттерне
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)

// Получить path value (Go 1.22+):
func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "user id: %s", id)
}
```

До Go 1.22 (и для более сложной маршрутизации используют chi, gorilla/mux):

```go
mux.HandleFunc("/users/", usersHandler)  // префикс-маршрут

func usersHandler(w http.ResponseWriter, r *http.Request) {
    // разбираем путь вручную
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) == 3 && parts[2] != "" {
        // /users/{id}
        id := parts[2]
        switch r.Method {
        case http.MethodGet:
            getUser(w, r, id)
        }
        return
    }
    http.NotFound(w, r)
}
```

---

## http.ResponseWriter — формирование ответа

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 1. Установить заголовки ПЕРЕД WriteHeader
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", generateID())

    // 2. Установить статус (по умолчанию 200)
    w.WriteHeader(http.StatusCreated)  // 201

    // 3. Записать тело
    json.NewEncoder(w).Encode(map[string]any{"id": 42})
}
```

Вспомогательные функции:

```go
// Ошибки:
http.Error(w, "Not Found", http.StatusNotFound)
http.NotFound(w, r)

// Редирект:
http.Redirect(w, r, "/new-path", http.StatusMovedPermanently)

// Отдача файла:
http.ServeFile(w, r, "index.html")
```

---

## Middleware — цепочка обработчиков

Middleware оборачивает `Handler` добавляя логику до/после:

```go
type Middleware func(http.Handler) http.Handler

// Логирование:
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("→ %s %s", r.Method, r.URL.Path)

        next.ServeHTTP(w, r)

        log.Printf("← %s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

// Аутентификация:
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Recovery от паники:
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### Цепочка middleware

```go
func chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /users/{id}", getUser)

    handler := chain(mux,
        recoveryMiddleware,
        loggingMiddleware,
        authMiddleware,
    )

    http.ListenAndServe(":8080", handler)
}
```

---

## http.Server с настройками

`http.ListenAndServe` использует дефолтный сервер без таймаутов — опасно в production:

```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  5 * time.Second,   // время на чтение тела запроса
    WriteTimeout: 10 * time.Second,  // время на запись ответа
    IdleTimeout:  120 * time.Second, // keep-alive соединения
    ReadHeaderTimeout: 2 * time.Second,
}

if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
    log.Fatal(err)
}
```

---

## Graceful Shutdown — мягкое завершение

Graceful shutdown: дождаться завершения активных запросов перед остановкой сервера:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(2 * time.Second)  // имитация долгой обработки
        fmt.Fprintln(w, "ok")
    })

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // Запускаем сервер в горутине
    go func() {
        log.Println("сервер запущен на :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("ошибка сервера: %v", err)
        }
    }()

    // Ждём сигнал остановки
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("получен сигнал остановки...")

    // Даём 30 секунд на завершение активных запросов
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("принудительное завершение: %v", err)
    }

    log.Println("сервер остановлен")
}
```

---

## Работа с JSON API

```go
// Декодирование запроса:
func createUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    // Валидация:
    if req.Name == "" || req.Email == "" {
        http.Error(w, "name and email required", http.StatusUnprocessableEntity)
        return
    }

    user := User{ID: 1, Name: req.Name, Email: req.Email}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// Вспомогательная функция:
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}
```

---

## Итог

- `http.Handler` — интерфейс с одним методом `ServeHTTP`; всё строится вокруг него
- `http.ServeMux` — встроенный маршрутизатор; Go 1.22+ поддерживает метод и path params
- Middleware — функции `func(Handler) Handler`; цепочкой добавляют логирование, auth, recovery
- `http.Server` с таймаутами — обязательно для production (без таймаутов = уязвимость)
- Graceful shutdown: `srv.Shutdown(ctx)` ждёт завершения активных запросов
- `json.NewDecoder(r.Body).Decode(&v)` — потоковое чтение JSON без буферизации в память
