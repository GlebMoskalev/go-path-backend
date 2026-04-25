---
title: "Кастомные ошибки"
description: "Struct-ошибки, методы Error() и Unwrap(), когда нужны кастомные ошибки"
order: 2
---

# Кастомные ошибки

Простых `errors.New` хватает для базовых случаев. Когда ошибка должна нести дополнительные данные — создавай struct-ошибки.

## Когда нужны кастомные ошибки

- Нужно передать структурированные данные (код ошибки, поле, ресурс)
- Вызывающий код должен различать типы ошибок через type assertion
- Нужно добавить методы (например, `IsTemporary() bool`)

---

## Struct-ошибки

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("ошибка валидации поля %q: %s", e.Field, e.Message)
}

// Создание:
func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{Field: "age", Message: "должно быть неотрицательным"}
    }
    if age > 150 {
        return &ValidationError{Field: "age", Message: "слишком большое значение"}
    }
    return nil
}

// Использование:
err := validateAge(-5)
if err != nil {
    fmt.Println(err)  // ошибка валидации поля "age": должно быть неотрицательным
}
```

---

## Извлечение данных через errors.As

`errors.As` позволяет извлечь конкретный тип ошибки из цепочки (включая wrapped):

```go
err := validateAge(-5)

var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Printf("Поле: %s, Сообщение: %s\n", ve.Field, ve.Message)
    // Поле: age, Сообщение: должно быть неотрицательным
}
```

---

## Метод Unwrap() — для вложенных ошибок

Если кастомная ошибка оборачивает другую:

```go
type DatabaseError struct {
    Op  string  // операция: "query", "insert", etc.
    Err error   // исходная ошибка
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("ошибка БД при %s: %v", e.Op, e.Err)
}

// Unwrap позволяет errors.Is/errors.As заглядывать внутрь:
func (e *DatabaseError) Unwrap() error {
    return e.Err
}

// Использование:
var ErrConnFailed = errors.New("соединение разорвано")

func queryUser(id int) (*User, error) {
    err := db.Query("SELECT * FROM users WHERE id = ?", id)
    if err != nil {
        return nil, &DatabaseError{Op: "query", Err: fmt.Errorf("%w: %w", ErrConnFailed, err)}
    }
    return &User{}, nil
}

err := queryUser(42)
if errors.Is(err, ErrConnFailed) {
    // errors.Is проходит цепочку через Unwrap()
    fmt.Println("Проблема с соединением")
}
```

---

## Паттерн: ошибки с кодом (HTTP-стиль)

```go
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) IsClientError() bool { return e.Code >= 400 && e.Code < 500 }
func (e *AppError) IsServerError() bool { return e.Code >= 500 }

// Конструкторы:
func NotFoundError(resource string) *AppError {
    return &AppError{Code: 404, Message: resource + " не найден"}
}

func InternalError(err error) *AppError {
    return &AppError{Code: 500, Message: "внутренняя ошибка", Err: err}
}
```

---

## Несколько ошибок: multierror

Иногда нужно накопить несколько ошибок:

```go
type MultiError struct {
    Errors []error
}

func (m *MultiError) Error() string {
    msgs := make([]string, len(m.Errors))
    for i, err := range m.Errors {
        msgs[i] = err.Error()
    }
    return strings.Join(msgs, "; ")
}

func (m *MultiError) Add(err error) {
    if err != nil {
        m.Errors = append(m.Errors, err)
    }
}

func (m *MultiError) ToError() error {
    if len(m.Errors) == 0 {
        return nil
    }
    return m
}

func validateUser(u User) error {
    var errs MultiError
    if u.Name == "" {
        errs.Add(errors.New("имя обязательно"))
    }
    if u.Age < 0 {
        errs.Add(fmt.Errorf("возраст %d недопустим", u.Age))
    }
    if !isValidEmail(u.Email) {
        errs.Add(fmt.Errorf("неверный email: %s", u.Email))
    }
    return errs.ToError()
}
```

В Go 1.20+ появился `errors.Join`:

```go
import "errors"

err1 := errors.New("первая ошибка")
err2 := errors.New("вторая ошибка")
combined := errors.Join(err1, err2)
fmt.Println(combined)
// первая ошибка
// вторая ошибка

errors.Is(combined, err1)  // true
```

---

## Итог

- Struct-ошибки несут дополнительные данные и поддерживают type assertion
- Метод `Error() string` — обязателен (реализует интерфейс)
- Метод `Unwrap() error` — опционален; позволяет `errors.Is/As` проходить цепочки
- `errors.As(err, &target)` — извлечь конкретный тип из цепочки ошибок
- Используй кастомные ошибки когда вызывающий код должен знать детали ошибки
- `errors.Join` (Go 1.20+) — объединение нескольких ошибок
