---
title: "Обработка ошибок"
description: "Идиоматичная обработка ошибок: error, fmt.Errorf, errors.Is/As, panic/recover"
order: 2
---

# Обработка ошибок

В Go нет исключений. Ошибки — это обычные значения, реализующие интерфейс `error`.

## Интерфейс error

```go
type error interface {
    Error() string
}
```

## Возврат ошибки

```go
func openFile(name string) (*os.File, error) {
    f, err := os.Open(name)
    if err != nil {
        return nil, err
    }
    return f, nil
}
```

## Проверка ошибок

```go
f, err := openFile("data.txt")
if err != nil {
    log.Fatal(err)
}
defer f.Close()
```

> **Правило:** всегда проверяйте ошибки. Go-линтеры предупреждают о непроверенных ошибках.

## Создание ошибок

### errors.New

```go
import "errors"

var ErrNotFound = errors.New("не найдено")
```

### fmt.Errorf

```go
func findUser(id int) (*User, error) {
    // ...
    return nil, fmt.Errorf("пользователь %d не найден", id)
}
```

### Оборачивание (wrapping)

```go
func loadConfig(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("loadConfig: %w", err)  // %w — оборачивание
    }
    // ...
    return nil
}
```

## errors.Is — проверка по значению

```go
err := loadConfig("missing.yaml")

if errors.Is(err, os.ErrNotExist) {
    fmt.Println("файл не существует")
}
```

`errors.Is` разворачивает всю цепочку обёрток.

## errors.As — проверка по типу

```go
var pathErr *os.PathError

if errors.As(err, &pathErr) {
    fmt.Println("путь:", pathErr.Path)
}
```

## Собственный тип ошибки

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("поле %s: %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{
            Field:   "age",
            Message: "не может быть отрицательным",
        }
    }
    return nil
}
```

## Sentinel-ошибки

```go
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)

func getUser(id int) (*User, error) {
    // ...
    return nil, ErrNotFound
}

// Использование:
if errors.Is(err, ErrNotFound) {
    // обработка
}
```

## panic и recover

`panic` — аварийная остановка. Используется для невосстановимых ошибок.

```go
func mustParseURL(raw string) *url.URL {
    u, err := url.Parse(raw)
    if err != nil {
        panic(err)
    }
    return u
}
```

`recover` — перехват panic внутри `defer`:

```go
func safeDiv(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    return a / b, nil
}
```

> **Правило:** не используйте `panic` для обычной обработки ошибок. Только для действительно невосстановимых ситуаций (нарушение инвариантов, критические ошибки инициализации).
