---
title: "context.Context"
description: "Background/TODO, WithCancel/WithTimeout/WithDeadline/WithValue, распространение отмены"
order: 2
---

# context.Context

`context.Context` — стандартный способ передавать сигналы отмены, дедлайны и значения через цепочку вызовов и горутин.

## Зачем нужен context

```go
// БЕЗ context: нет способа отменить операцию
func fetchData(url string) ([]byte, error) {
    resp, err := http.Get(url)
    // что если пользователь нажал «Отмена»? запрос продолжится
    // ...
}

// С context: отмена распространяется автоматически
func fetchData(ctx context.Context, url string) ([]byte, error) {
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    resp, err := http.DefaultClient.Do(req)
    // если ctx отменён — http.Do вернёт ошибку немедленно
    // ...
}
```

---

## Создание корневого контекста

```go
import "context"

// Background: корневой контекст для main, серверов, глобальных операций
ctx := context.Background()

// TODO: заглушка когда контекст ещё не определён; помечает незаконченный код
ctx := context.TODO()
```

**Никогда не передавай `nil` в качестве context** — используй `context.Background()` или `context.TODO()`.

---

## WithCancel — ручная отмена

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    // имитируем отмену через 2 секунды
    time.Sleep(2 * time.Second)
    cancel()  // отменяем контекст
}()

select {
case <-time.After(5 * time.Second):
    fmt.Println("работа завершена")
case <-ctx.Done():
    fmt.Println("отменено:", ctx.Err())  // context canceled
}
```

`cancel()` следует вызывать всегда, даже если работа завершилась успешно — это освобождает ресурсы. Обычно через `defer cancel()`:

```go
func processRequest(parentCtx context.Context) error {
    ctx, cancel := context.WithCancel(parentCtx)
    defer cancel()  // гарантированная очистка

    // ...
    return nil
}
```

---

## WithTimeout и WithDeadline

```go
// Отмена через 5 секунд
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Отмена в конкретный момент времени
deadline := time.Now().Add(5 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()
```

Пример с HTTP-запросом:

```go
func fetchWithTimeout(url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        // err содержит информацию об отмене если ctx истёк
        return nil, fmt.Errorf("запрос не выполнен: %w", err)
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### Проверка причины отмены

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

select {
case <-ctx.Done():
    switch ctx.Err() {
    case context.DeadlineExceeded:
        fmt.Println("дедлайн истёк")
    case context.Canceled:
        fmt.Println("отменено вручную")
    }
}
```

---

## Распространение отмены

Контексты образуют дерево: отмена родителя отменяет всех потомков:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Дочерний контекст — с более коротким таймаутом
subCtx, subCancel := context.WithTimeout(ctx, 2*time.Second)
defer subCancel()

// subCtx отменится через 2 секунды ИЛИ когда ctx отменится (через 10с)
// — что наступит раньше
```

Практический пример — HTTP-сервер:

```go
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // r.Context() уже содержит отмену при разрыве соединения клиента
    ctx := r.Context()

    result, err := h.service.Process(ctx, r.URL.Query().Get("id"))
    if err != nil {
        if errors.Is(err, context.Canceled) {
            return  // клиент ушёл, ничего не делаем
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(result)
}
```

---

## WithValue — передача значений

`WithValue` добавляет пару ключ-значение в контекст:

```go
// Ключ должен быть unexported типом — чтобы избежать коллизий между пакетами
type contextKey string

const (
    keyUserID    contextKey = "userID"
    keyRequestID contextKey = "requestID"
)

func withUserID(ctx context.Context, id int) context.Context {
    return context.WithValue(ctx, keyUserID, id)
}

func userIDFromContext(ctx context.Context) (int, bool) {
    id, ok := ctx.Value(keyUserID).(int)
    return id, ok
}
```

Использование в middleware:

```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        userID, err := validateToken(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Добавляем userID в контекст — доступен всем вниз по цепочке
        ctx := withUserID(r.Context(), userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    userID, ok := userIDFromContext(r.Context())
    if !ok {
        http.Error(w, "user not found in context", http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "hello user %d", userID)
}
```

### Правила использования WithValue

- Только для данных уровня запроса: request ID, user ID, трассировка
- Не для обязательных параметров функций — передавай их явно
- Тип ключа — всегда unexported (иначе риск коллизии с другими пакетами)
- Значение должно быть безопасно для конкурентного чтения (контекст может читаться из нескольких горутин)

---

## context в горутинах

```go
func processItems(ctx context.Context, items []Item) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(items))

    for _, item := range items {
        wg.Add(1)
        go func(item Item) {
            defer wg.Done()
            if err := processItem(ctx, item); err != nil {
                errCh <- err
            }
        }(item)
    }

    wg.Wait()
    close(errCh)

    for err := range errCh {
        if err != nil {
            return err
        }
    }
    return nil
}

func processItem(ctx context.Context, item Item) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // выполняем работу
    return doWork(item)
}
```

---

## Итог

- `context.Background()` — корневой контекст; `context.TODO()` — заглушка
- `WithCancel` — ручная отмена через `cancel()`; всегда `defer cancel()`
- `WithTimeout(ctx, d)` — отмена через длительность d
- `WithDeadline(ctx, t)` — отмена в момент времени t
- Отмена родителя автоматически отменяет всех потомков
- `ctx.Done()` — канал, закрывается при отмене; `ctx.Err()` — причина
- `WithValue` — для данных уровня запроса; ключ — unexported тип
- Первый параметр функции с I/O — всегда `ctx context.Context`
