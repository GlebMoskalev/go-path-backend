---
title: "switch"
description: "switch без условия, fallthrough, type switch, сравнение с if-else chains"
order: 2
---

# switch

`switch` в Go мощнее, чем в C/Java: нет автоматического fallthrough, можно использовать без выражения, поддерживает несколько значений в одном case и работу с типами.

## Базовый синтаксис

```go
switch выражение {
case значение1:
    // ...
case значение2:
    // ...
default:
    // если ни один case не совпал
}
```

**Ключевое отличие от C/Java**: в Go `break` в конце каждого case **по умолчанию подразумевается**. Выполнив один case, `switch` завершается — не переходит в следующий:

```go
x := 2

switch x {
case 1:
    fmt.Println("один")
case 2:
    fmt.Println("два")  // выводит только это
case 3:
    fmt.Println("три")  // это не выполняется
}
```

### Несколько значений в одном case

```go
func classify(day string) string {
    switch day {
    case "Понедельник", "Вторник", "Среда", "Четверг", "Пятница":
        return "будний день"
    case "Суббота", "Воскресенье":
        return "выходной"
    default:
        return "неизвестный день"
    }
}
```

---

## switch с инициализирующим оператором

Аналогично `if`, `switch` поддерживает init statement:

```go
switch os := runtime.GOOS; os {
case "linux":
    fmt.Println("Linux")
case "darwin":
    fmt.Println("macOS")
default:
    fmt.Printf("Другая ОС: %s\n", os)
}
```

---

## switch без выражения

`switch` без условного выражения — более читаемая замена длинным цепочкам `if/else if`:

```go
func httpStatus(code int) string {
    switch {
    case code >= 500:
        return "ошибка сервера"
    case code >= 400:
        return "ошибка клиента"
    case code >= 300:
        return "перенаправление"
    case code >= 200:
        return "успех"
    default:
        return "неизвестный статус"
    }
}
```

Каждый `case` — произвольное булево выражение. Проверяются сверху вниз, выполняется первый истинный.

Сравни с `if/else if` — `switch {}` читается лучше при 3+ ветвях:

```go
// Менее читаемо:
if code >= 500 {
    return "ошибка сервера"
} else if code >= 400 {
    return "ошибка клиента"
} else if code >= 300 {
    return "перенаправление"
} else if code >= 200 {
    return "успех"
} else {
    return "неизвестный статус"
}
```

---

## fallthrough — явный переход

Когда нужно поведение C-style (выполнить следующий case), используй `fallthrough`:

```go
n := 2

switch n {
case 1:
    fmt.Println("один")
    fallthrough
case 2:
    fmt.Println("два или продолжение от одного")
    fallthrough
case 3:
    fmt.Println("три или продолжение")
case 4:
    fmt.Println("четыре")
}
// Если n == 2:
// два или продолжение от одного
// три или продолжение
// (case 4 не выполняется — нет fallthrough после case 3)
```

**Важно**: `fallthrough` выполняет следующий case **безусловно**, не проверяя его значение. В последнем case блока использовать нельзя — компилятор выдаст ошибку. Используется редко — только когда действительно нужна такая логика.

---

## Type switch

Позволяет переключаться по **типу** значения интерфейса:

```go
func describe(i interface{}) string {
    switch v := i.(type) {
    case int:
        return fmt.Sprintf("int: %d", v)
    case string:
        return fmt.Sprintf("string: %q (длина %d)", v, len(v))
    case bool:
        if v {
            return "bool: true"
        }
        return "bool: false"
    case []int:
        return fmt.Sprintf("[]int длиной %d", len(v))
    case nil:
        return "nil"
    default:
        return fmt.Sprintf("неизвестный тип: %T", v)
    }
}

func main() {
    fmt.Println(describe(42))        // int: 42
    fmt.Println(describe("hello"))   // string: "hello" (длина 5)
    fmt.Println(describe(true))      // bool: true
    fmt.Println(describe([]int{1,2})) // []int длиной 2
    fmt.Println(describe(nil))       // nil
}
```

В каждом `case` переменная `v` имеет конкретный тип — это делает код безопасным без лишних приведений.

> 💡 Интерфейсы и type assertion подробно разберём в главе 7 «Интерфейсы».

---

## Полный практический пример

```go
package main

import (
    "fmt"
    "time"
)

type Weekday = time.Weekday

func schedule(day Weekday, hour int) string {
    switch day {
    case time.Saturday, time.Sunday:
        return "выходной, можно отдыхать"
    }

    switch {
    case hour < 9:
        return "слишком рано, офис ещё закрыт"
    case hour < 12:
        return "утренняя сессия работы"
    case hour < 14:
        return "обеденный перерыв"
    case hour < 18:
        return "послеобеденная сессия"
    default:
        return "рабочий день окончен"
    }
}

func main() {
    now := time.Now()
    fmt.Println(schedule(now.Weekday(), now.Hour()))
}
```

---

## Итог

- `switch` в Go не переходит в следующий case — `break` не нужен
- Несколько значений в одном case: `case "a", "b", "c":`
- `switch {}` без выражения = удобная замена длинному `if/else if`
- `fallthrough` — явный переход в следующий case, используется редко
- Type switch — переключение по типу: `switch v := i.(type)`
- `switch` с init statement: `switch x := expr; x {`
