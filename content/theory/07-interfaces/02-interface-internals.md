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

Когда ты вызываешь метод напрямую — компилятор знает точный тип и сразу подставляет адрес нужной функции. Это называется статический вызов.

```go
c := Circle{Radius: 5}
c.Area()  // компилятор знает: это Circle.Area — вызов прямой
```

Когда ты вызываешь метод через интерфейс — компилятор не знает заранее, что лежит внутри: `Circle`, `Rectangle` или что-то ещё. Поэтому при каждом вызове Go сначала смотрит в itab (таблицу, где хранятся указатели на методы конкретного типа), находит нужную функцию и только потом вызывает. Это дополнительный шаг — косвенный вызов.

```go
var s Shape = Circle{Radius: 5}
s.Area()  // Go сначала смотрит в itab: «для Circle метод Area вот здесь» → вызов
```

На практике разница в скорости между прямым и косвенным вызовом — наносекунды, и в большинстве программ это вообще не заметно. Беспокоиться об этом стоит только если profiler показывает, что именно этот вызов — узкое место.

---

## Интерфейсы и nil: ловушка типизированного nil

Это одна из самых частых неожиданностей в Go. Вспомни как устроен интерфейс — пара `(type, value)`. Интерфейс считается nil только если **оба** поля пустые. Если тип есть, а значение nil — интерфейс всё равно не nil.

Посмотри на этот код:

```go
func getUser(id int) (*User, error) {
    var err *ValidationError  // err = (*ValidationError)(nil) — тип есть, значение nil
    if id <= 0 {
        err = &ValidationError{Message: "invalid id"}
    }
    return user, err
}
```

Когда `id > 0`, мы возвращаем `err` — а это `(*ValidationError)(nil)`. При присвоении в интерфейс `error` Go запишет: `type = *ValidationError, value = nil`. И проверка у вызывающего сломается:

```go
user, err := getUser(42)
if err != nil {  // TRUE! хотя ошибки нет
    log.Fatal(err)  // программа падает без причины
}
```

`err != nil` вернёт `true`, потому что тип `*ValidationError` есть — интерфейс не пустой. Хотя ошибки никакой нет.

Правильно — возвращать явный `nil` без типа:

```go
func getUser(id int) (*User, error) {
    if id <= 0 {
        return nil, &ValidationError{Message: "invalid id"}
    }
    return user, nil  // nil без типа — интерфейс будет (nil, nil) → настоящий nil
}
```

Здесь `nil` — это не типизированный nil, а просто пустой интерфейс: `(type=nil, value=nil)`. Проверка `err != nil` вернёт `false` — всё корректно.

То же самое касается любых возвращаемых интерфейсов, не только `error`:

```go
// Плохо: возвращает (*os.File)(nil) — Writer не nil, но файла нет
func newWriter(path string) io.Writer {
    var w *os.File
    if path != "" {
        w, _ = os.Create(path)
    }
    return w
}

// Хорошо: явный nil — Writer настоящий nil
func newWriter(path string) io.Writer {
    if path == "" {
        return nil
    }
    w, _ := os.Create(path)
    return w
}
```

**Правило**: если функция возвращает интерфейс — возвращай `nil` напрямую, никогда не возвращай переменную конкретного типа с нулевым значением.

---

## Итог

- Интерфейс = пара `(type, value)`
- Интерфейс == nil только когда **оба** компонента nil
- Если в интерфейс помещён nil-указатель конкретного типа — интерфейс **не nil**
- Никогда не возвращай конкретный nil-указатель как интерфейс — возвращай `nil` напрямую
- Динамическая диспетчеризация — небольшой overhead, обычно несущественный
