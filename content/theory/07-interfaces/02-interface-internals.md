---
title: "Внутреннее устройство интерфейсов"
description: "Пара (тип, значение), nil interface vs interface с nil значением — классическая ловушка"
order: 2
---

# Внутреннее устройство интерфейсов

Знание внутреннего устройства интерфейсов необходимо, чтобы понять одну из самых коварных ловушек Go — разницу между nil interface и interface с nil значением.

## Пара (type, value)

Значение интерфейса — это пара из двух компонентов:
1. **Динамический тип** (type) — конкретный тип хранимого значения
2. **Динамическое значение** (value) — само значение

```
Интерфейс:
┌───────────────────┐
│      *type        │  → указатель на информацию о типе
├───────────────────┤
│      *data        │  → указатель на данные
└───────────────────┘
```

```go
var r io.Reader

// r = (type: nil, value: nil) — оба nil, интерфейс nil
fmt.Println(r == nil)  // true

f, _ := os.Open("file.txt")
r = f
// r = (type: *os.File, value: <pointer to file>)
fmt.Println(r == nil)  // false
```

---

## Классическая ловушка: nil в интерфейсе

Это одна из самых частых ошибок в Go-коде. Интерфейс nil **только тогда**, когда **оба** компонента nil.

Если в интерфейс помещён nil-указатель конкретного типа — интерфейс не nil:

```go
package main

import "fmt"

type MyError struct {
    Code int
}

func (e *MyError) Error() string {
    return fmt.Sprintf("ошибка %d", e.Code)
}

func getError(fail bool) error {
    var err *MyError  // nil указатель на MyError
    if fail {
        err = &MyError{Code: 404}
    }
    return err  // ОШИБКА: возвращаем (*MyError)(nil) как error
}

func main() {
    err := getError(false)
    if err != nil {
        fmt.Println("Ошибка:", err)  // ЭТО ВЫПОЛНИТСЯ!
        // error = (type: *MyError, value: nil) — не nil!
    }
}
// Вывод: Ошибка: <nil>
```

**Что происходит:**
- `var err *MyError` — nil-указатель на MyError
- `return err` — помещаем `(*MyError)(nil)` в интерфейс `error`
- Интерфейс = `(type: *MyError, value: nil)` — тип НЕ nil!
- `err != nil` — true, хотя реального значения нет

### Правильное решение

```go
func getError(fail bool) error {
    if fail {
        return &MyError{Code: 404}
    }
    return nil  // явно возвращаем nil интерфейс, а не nil-указатель
}

func main() {
    err := getError(false)
    if err != nil {
        fmt.Println("Ошибка:", err)  // НЕ выполнится
    }
    fmt.Println("Успех")  // Успех
}
```

**Правило:** никогда не возвращай конкретный nil-указатель там, где ожидается интерфейс. Возвращай чистый `nil`.

---

## Визуализация разных состояний

```go
// 1. Nil интерфейс (оба компонента nil):
var e error
// e = (type: nil, value: nil)
// e == nil → true

// 2. Интерфейс с nil-значением конкретного типа:
var p *MyError = nil
e = p
// e = (type: *MyError, value: nil)
// e == nil → false !!!

// 3. Интерфейс с реальным значением:
e = &MyError{404}
// e = (type: *MyError, value: 0xc000...pointing to MyError{404})
// e == nil → false
```

---

## Как правильно сравнивать

Если нужно проверить, является ли значение в интерфейсе nil, используй type assertion:

```go
func isNilValue(i interface{}) bool {
    if i == nil {
        return true  // nil интерфейс
    }
    
    v := reflect.ValueOf(i)
    switch v.Kind() {
    case reflect.Ptr, reflect.Chan, reflect.Func,
         reflect.Interface, reflect.Map, reflect.Slice:
        return v.IsNil()
    }
    return false
}
```

Но на практике лучше просто не допускать ситуации, где это нужно — следуй правилу выше.

---

## Динамическая диспетчеризация

Когда вызывается метод через интерфейс, Go делает косвенный вызов через таблицу методов (vtable-подобный механизм):

```go
var s Shape = Circle{Radius: 5}
s.Area()  // косвенный вызов через таблицу методов
```

Это чуть медленнее прямого вызова метода, но разница обычно незначительна. Для горячих путей с критичной производительностью — делай benchmark.

---

## Интерфейсы и nil: практические рекомендации

```go
// НЕ ДЕЛАЙ ТАК:
func getUser(id int) (*User, error) {
    var err *ValidationError
    if id <= 0 {
        err = &ValidationError{Message: "invalid id"}
    }
    // ...
    return user, err  // err — (*ValidationError)(nil) если id > 0
}

// ДЕЛАЙ ТАК:
func getUser(id int) (*User, error) {
    if id <= 0 {
        return nil, &ValidationError{Message: "invalid id"}
    }
    // ...
    return user, nil  // явный nil
}
```

```go
// НЕ ДЕЛАЙ ТАК (в функциях, возвращающих интерфейс):
func newWriter(path string) io.Writer {
    var w *os.File
    if path != "" {
        w, _ = os.Create(path)
    }
    return w  // возвращает (*os.File)(nil) — не nil Writer!
}

// ДЕЛАЙ ТАК:
func newWriter(path string) io.Writer {
    if path == "" {
        return nil  // или io.Discard
    }
    w, _ := os.Create(path)
    return w
}
```

---

## Итог

- Интерфейс = пара `(type, value)`
- Интерфейс == nil только когда **оба** компонента nil
- Если в интерфейс помещён nil-указатель конкретного типа — интерфейс **не nil**
- Никогда не возвращай конкретный nil-указатель как интерфейс — возвращай `nil` напрямую
- Динамическая диспетчеризация — небольшой overhead, обычно несущественный
