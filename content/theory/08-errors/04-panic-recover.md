---
title: "panic и recover"
description: "panic vs error, recover() в defer, когда panic уместен, не использовать для control flow"
order: 4
---

# panic и recover

`panic` — механизм для ситуаций, которые не должны происходить в корректно работающей программе. Это не замена обработке ошибок.

## Что такое panic

`panic` немедленно прекращает выполнение текущей функции, раскручивает стек (выполняя все `defer`) и завершает программу, если паника не перехвачена `recover`.

```go
func main() {
    fmt.Println("до паники")
    panic("что-то пошло не так")
    fmt.Println("это не выполнится")
}
// до паники
// panic: что-то пошло не так
//
// goroutine 1 [running]:
// main.main()
//         /tmp/main.go:5 +0x65
// exit status 2
```

Go сам вызывает `panic` при:
- Обращении по nil-указателю
- Выходе за границы слайса/массива
- Делении на ноль (для целых)
- Неверном type assertion без comma-ok

```go
var p *int
*p = 5  // panic: runtime error: invalid memory address or nil pointer dereference

s := []int{1, 2, 3}
s[10]   // panic: runtime error: index out of range [10] with length 3

var i interface{} = "hello"
n := i.(int)  // panic: interface conversion: interface {} is string, not int
```

---

## Когда использовать panic

**Использовать panic** уместно только для:

1. **Невозможных состояний** — нарушение инвариантов программы, которое говорит о баге:
```go
func divide(a, b int) int {
    if b == 0 {
        panic("divide: делитель не может быть нулём — это баг вызывающего кода")
    }
    return a / b
}
```

2. **Ошибок инициализации** — если программа не может продолжить работу:
```go
func mustParseTemplate(text string) *template.Template {
    t, err := template.New("").Parse(text)
    if err != nil {
        panic(fmt.Sprintf("mustParseTemplate: неверный шаблон: %v", err))
    }
    return t
}

// Используется при инициализации пакета:
var homeTmpl = mustParseTemplate(`<html>{{.Name}}</html>`)
```

3. **Несовместимых версий** / нереализованного кода:
```go
func unimplemented() {
    panic("not implemented")
}
```

**НЕ используй panic** для:
- Обычных ошибок (нет файла, неверный ввод, ошибка сети)
- Control flow (вместо `return error`)
- Обработки ошибок внешних API

---

## recover — перехват паники

`recover()` останавливает раскрутку стека и возвращает значение, переданное в `panic`. Работает **только внутри defer**:

```go
func safeDiv(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("паника: %v", r)
        }
    }()

    result = a / b  // паника при b == 0
    return
}

result, err := safeDiv(10, 0)
fmt.Println(result, err)  // 0 паника: runtime error: integer divide by zero

result, err = safeDiv(10, 2)
fmt.Println(result, err)  // 5 <nil>
```

### recover не работает в другой горутине

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("восстановились:", r)  // НЕ выполнится
        }
    }()

    go func() {
        panic("паника в горутине")  // завершит программу!
    }()

    time.Sleep(time.Second)
}
// panic: паника в горутине
```

Каждая горутина должна сама обрабатывать свои паники.

---

## Паттерн: защита HTTP-обработчика

Стандартный паттерн в веб-серверах:

```go
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Логируем трейс:
                buf := make([]byte, 4096)
                n := runtime.Stack(buf, false)
                log.Printf("PANIC: %v\n%s", err, buf[:n])

                http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

Это позволяет серверу продолжить работу после паники в одном обработчике.

---

## Паттерн: panic/recover как внутренний механизм

Иногда `panic`/`recover` используют внутри пакета как удобный способ выхода из глубокой рекурсии, но наружу всегда возвращают `error`:

```go
// Только для internal use:
type parseError struct{ msg string }

func (p *parser) parse() (result AST, err error) {
    defer func() {
        if r := recover(); r != nil {
            if pe, ok := r.(parseError); ok {
                err = errors.New(pe.msg)
            } else {
                panic(r)  // не наша паника — пробрасываем дальше
            }
        }
    }()

    return p.parseExpr(), nil
}

func (p *parser) expect(token string) {
    if p.next() != token {
        panic(parseError{fmt.Sprintf("ожидался %q", token)})
    }
}
```

Это допустимо, но только **внутри пакета**. Публичный API должен возвращать `error`.

---

## Различие: panic vs log.Fatal

```go
// panic — раскручивает стек, выполняет defer:
panic("что-то сломалось")

// log.Fatal — выводит сообщение и вызывает os.Exit(1):
// defer НЕ выполняются!
log.Fatal("фатальная ошибка")

// os.Exit — немедленное завершение, defer НЕ выполняются:
os.Exit(1)
```

Используй `log.Fatal` / `os.Exit` только в `main()` при неустранимой ошибке запуска.

---

## Итог

- `panic` — для программных ошибок (нарушение инвариантов, невозможные состояния)
- Не используй `panic` для обычной обработки ошибок
- `recover()` работает только внутри `defer`
- `recover()` не перехватывает паники из других горутин
- Пробрасывай не свои паники: `if _, ok := r.(myPanicType); !ok { panic(r) }`
- В HTTP-серверах обязательно используй recovery middleware
- `log.Fatal` и `os.Exit` не выполняют `defer` — используй только в main
