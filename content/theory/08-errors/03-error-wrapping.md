---
title: "Wrapping ошибок"
description: "fmt.Errorf с %w, errors.Is, errors.As, цепочки ошибок (Go 1.13+)"
order: 3
---

# Wrapping ошибок

Go 1.13 (2019) добавил механизм оборачивания ошибок — стандартный способ добавлять контекст к ошибкам без потери исходной информации.

## Зачем оборачивать ошибки

Без оборачивания теряется контекст:

```go
// Плохо — потеря контекста:
func processOrder(id int) error {
    _, err := db.GetOrder(id)
    if err != nil {
        return err  // какой заказ? откуда ошибка?
    }
    return nil
}
```

С оборачиванием — полная цепочка контекста:

```go
// Хорошо — контекст сохраняется:
func processOrder(id int) error {
    _, err := db.GetOrder(id)
    if err != nil {
        return fmt.Errorf("processOrder(id=%d): %w", id, err)
    }
    return nil
}
// Ошибка: processOrder(id=42): db: query failed: connection refused
```

---

## fmt.Errorf с %w

Оператор `%w` (wrap) создаёт обёрнутую ошибку:

```go
originalErr := errors.New("соединение разорвано")
wrappedErr := fmt.Errorf("не удалось подключиться к БД: %w", originalErr)

fmt.Println(wrappedErr)
// не удалось подключиться к БД: соединение разорвано

// wrappedErr содержит originalErr внутри:
fmt.Println(errors.Is(wrappedErr, originalErr))  // true
```

Разница между `%v` и `%w`:

```go
err := errors.New("исходная ошибка")

withV := fmt.Errorf("контекст: %v", err)  // просто строка
withW := fmt.Errorf("контекст: %w", err)  // обёрнутая ошибка

fmt.Println(errors.Is(withV, err))  // false! потеряна связь
fmt.Println(errors.Is(withW, err))  // true — связь сохранена
```

---

## errors.Is — проверка по цепочке

`errors.Is(err, target)` проверяет, содержит ли цепочка ошибок `target`:

```go
var ErrNotFound = errors.New("не найдено")

func getUser(id int) (*User, error) {
    return nil, fmt.Errorf("getUser: %w", ErrNotFound)
}

func processRequest(id int) error {
    _, err := getUser(id)
    if err != nil {
        return fmt.Errorf("processRequest: %w", err)
    }
    return nil
}

err := processRequest(99)
fmt.Println(err)
// processRequest: getUser: не найдено

// errors.Is раскручивает цепочку через Unwrap():
fmt.Println(errors.Is(err, ErrNotFound))  // true — несмотря на 2 уровня обёртки
```

`errors.Is` последовательно вызывает `Unwrap()` на каждой ошибке в цепочке, пока не найдёт совпадение.

### Кастомный errors.Is

Если нужно нестандартное сравнение — реализуй метод `Is`:

```go
type HTTPError struct {
    Code int
}

func (e *HTTPError) Error() string { return fmt.Sprintf("HTTP %d", e.Code) }

// Кастомное сравнение: HTTPError равна другой HTTPError с тем же кодом
func (e *HTTPError) Is(target error) bool {
    t, ok := target.(*HTTPError)
    if !ok {
        return false
    }
    return e.Code == t.Code
}

err404 := &HTTPError{Code: 404}
wrappedErr := fmt.Errorf("page not found: %w", err404)

// Работает потому что HTTPError реализует Is:
errors.Is(wrappedErr, &HTTPError{Code: 404})  // true
errors.Is(wrappedErr, &HTTPError{Code: 500})  // false
```

---

## errors.As — извлечение типа из цепочки

`errors.As(err, &target)` ищет в цепочке ошибку, которую можно присвоить `target`:

```go
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s — %s", e.Field, e.Message)
}

func validate(name string) error {
    if name == "" {
        return &ValidationError{Field: "name", Message: "обязательное поле"}
    }
    return nil
}

func createUser(name string) error {
    if err := validate(name); err != nil {
        return fmt.Errorf("createUser: %w", err)  // оборачиваем
    }
    return nil
}

err := createUser("")

// errors.As раскручивает цепочку и присваивает в target:
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Printf("Поле %q: %s\n", ve.Field, ve.Message)
    // Поле "name": обязательное поле
}
```

**Разница между `errors.Is` и `errors.As`:**
- `errors.Is` — проверяет равенство конкретному экземпляру (sentinel errors)
- `errors.As` — извлекает ошибку определённого **типа** (для struct-ошибок с данными)

---

## Практический пример: многоуровневая система

```go
package main

import (
    "errors"
    "fmt"
)

// Уровни:
// HTTP Handler → Service → Repository → Database

var ErrNotFound = errors.New("запись не найдена")

type DBError struct {
    Query string
    Err   error
}
func (e *DBError) Error() string { return fmt.Sprintf("db(%q): %v", e.Query, e.Err) }
func (e *DBError) Unwrap() error { return e.Err }

func dbQuery(q string) error {
    return &DBError{Query: q, Err: ErrNotFound}
}

func repoGetUser(id int) (*User, error) {
    err := dbQuery(fmt.Sprintf("SELECT * FROM users WHERE id=%d", id))
    if err != nil {
        return nil, fmt.Errorf("repo.GetUser(%d): %w", id, err)
    }
    return &User{}, nil
}

func serviceGetUser(id int) (*User, error) {
    user, err := repoGetUser(id)
    if err != nil {
        return nil, fmt.Errorf("service.GetUser: %w", err)
    }
    return user, nil
}

func main() {
    _, err := serviceGetUser(42)
    if err != nil {
        fmt.Println("Ошибка:", err)
        // service.GetUser: repo.GetUser(42): db("SELECT * FROM users WHERE id=42"): запись не найдена

        // Проверка sentinel:
        if errors.Is(err, ErrNotFound) {
            fmt.Println("→ Ресурс не найден")
        }

        // Извлечение данных БД:
        var dbErr *DBError
        if errors.As(err, &dbErr) {
            fmt.Println("→ Запрос БД:", dbErr.Query)
        }
    }
}
```

---

## Соглашения об оборачивании

**Добавляй контекст при подъёме по стеку:**
```go
// В каждой функции добавляй операцию/параметры:
return fmt.Errorf("service.ProcessOrder(id=%d): %w", id, err)
return fmt.Errorf("repo.GetUser(email=%q): %w", email, err)
```

**Не оборачивай дважды одно и то же:**
```go
// Плохо:
err := doSomething()
return fmt.Errorf("doSomething failed: %w", fmt.Errorf("operation error: %w", err))

// Хорошо — один слой контекста:
return fmt.Errorf("doSomething: %w", err)
```

---

## Итог

- `fmt.Errorf("контекст: %w", err)` — оборачивает ошибку, сохраняя цепочку
- `%v` — просто форматирует строку, связь с оригиналом теряется
- `errors.Is(err, target)` — проверяет цепочку на равенство sentinel-ошибке
- `errors.As(err, &target)` — извлекает ошибку конкретного типа из цепочки
- Добавляй контекст на каждом уровне: имя функции, параметры
- Go 1.20+: `errors.Join` для объединения нескольких ошибок
