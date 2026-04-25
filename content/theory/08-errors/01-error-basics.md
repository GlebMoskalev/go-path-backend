---
title: "Основы обработки ошибок"
description: "Тип error, errors.New, fmt.Errorf, идиома if err != nil, sentinel errors"
order: 1
---

# Основы обработки ошибок

Подход Go к ошибкам кардинально отличается от исключений в Java/Python/C++. Ошибки — обычные значения, которые возвращаются из функций и проверяются явно.

## Тип error

```go
type error interface {
    Error() string
}
```

`error` — встроенный интерфейс из одного метода. Любой тип с методом `Error() string` является ошибкой.

Нулевое значение интерфейса — `nil`. Если функция вернула `nil` как `error` — ошибки нет.

---

## Создание ошибок

### errors.New — простая строковая ошибка

```go
import "errors"

var ErrDivByZero = errors.New("деление на ноль")

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, ErrDivByZero
    }
    return a / b, nil
}
```

### fmt.Errorf — ошибка с форматированием

```go
func getUser(id int) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("недопустимый id пользователя: %d", id)
    }
    // ...
}
```

---

## Идиома if err != nil

Стандартный паттерн Go — проверять ошибку сразу после вызова функции:

```go
user, err := getUser(42)
if err != nil {
    return fmt.Errorf("не удалось получить пользователя: %w", err)
}

order, err := getOrder(user)
if err != nil {
    return fmt.Errorf("не удалось получить заказ: %w", err)
}

if err := saveOrder(order); err != nil {
    return fmt.Errorf("не удалось сохранить заказ: %w", err)
}
```

Да, это многословно. Но это намеренное решение: каждое место, где может возникнуть ошибка, явно обрабатывается.

**Почему это лучше исключений:**
- Ошибки — часть сигнатуры функции
- Нет скрытого потока управления
- Каждая ошибка обрабатывается в нужном контексте

---

## Sentinel errors — предопределённые ошибки

Sentinel errors — именованные переменные ошибок, которые используются для сравнения:

```go
package io

var EOF = errors.New("EOF")

// В пакете os:
var ErrNotExist = errors.New("file does not exist")
var ErrPermission = errors.New("permission denied")
```

Соглашение: sentinel errors объявляются на уровне пакета с префиксом `Err`:

```go
package store

import "errors"

var (
    ErrNotFound     = errors.New("запись не найдена")
    ErrAlreadyExists = errors.New("запись уже существует")
    ErrInvalidInput = errors.New("неверные входные данные")
)
```

### Сравнение с sentinel errors

```go
_, err := os.Open("nonexistent.txt")
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("файл не существует")
}
```

**Не используй `==` для сравнения ошибок** — используй `errors.Is` (он поддерживает wrapped errors):

```go
err := getUser(0)

// Плохо:
if err == ErrNotFound { ... }

// Хорошо:
if errors.Is(err, ErrNotFound) { ... }
```

---

## Практический пример: обработка ошибок в HTTP-обработчике

```go
package main

import (
    "errors"
    "fmt"
    "net/http"
)

var (
    ErrNotFound     = errors.New("не найдено")
    ErrUnauthorized = errors.New("не авторизован")
)

func getArticle(id int, userID int) (*Article, error) {
    if id <= 0 {
        return nil, fmt.Errorf("getArticle: неверный id %d", id)
    }
    // ... логика ...
    if userID == 0 {
        return nil, fmt.Errorf("getArticle: %w", ErrUnauthorized)
    }
    return &Article{}, nil
}

func handleArticle(w http.ResponseWriter, r *http.Request) {
    article, err := getArticle(1, 0)
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            http.Error(w, "Статья не найдена", http.StatusNotFound)
        case errors.Is(err, ErrUnauthorized):
            http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
        default:
            http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
        }
        return
    }
    fmt.Fprintf(w, "%+v", article)
}
```

---

## Типичные ошибки

**Ошибка 1**: Игнорировать ошибку через `_`.

```go
result, _ := riskyOperation()  // тихая ошибка!
// Иногда допустимо, но всегда должно быть обоснование
```

**Ошибка 2**: Сравнивать ошибки через `==` вместо `errors.Is`.

```go
// Может не сработать для wrapped errors:
if err == os.ErrNotExist { ... }

// Правильно:
if errors.Is(err, os.ErrNotExist) { ... }
```

**Ошибка 3**: Возвращать nil-указатель конкретного типа как error (рассмотрели в главе 7).

---

## Итог

- `error` — интерфейс с одним методом `Error() string`
- `errors.New("msg")` — простая ошибка; `fmt.Errorf("...: %w", err)` — с форматированием и wrapping
- Всегда проверяй `err != nil` после каждого вызова
- Sentinel errors — именованные ошибки пакетного уровня (`var ErrNotFound = errors.New(...)`)
- Используй `errors.Is` для сравнения, а не `==`
