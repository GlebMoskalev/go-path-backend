---
title: "Type Assertion"
description: "x.(T), comma-ok форма, паника при неверном утверждении"
order: 3
---

# Type Assertion

Type assertion позволяет извлечь конкретный тип из интерфейса. Это способ сказать: «я уверен, что в этом интерфейсе лежит значение типа T».

## Синтаксис

```go
// Форма без проверки — может вызвать панику:
value := iface.(T)

// Форма с проверкой (comma-ok) — безопасная:
value, ok := iface.(T)
```

---

## Простая форма — паника при ошибке

```go
var i interface{} = "hello"

s := i.(string)
fmt.Println(s)        // hello
fmt.Println(len(s))   // 5

// Если тип неверный — паника:
n := i.(int)
// panic: interface conversion: interface {} is string, not int
```

Используй только когда **уверен** в типе (например, сразу после type switch).

---

## Comma-ok — безопасная форма

```go
var i interface{} = "hello"

// Проверка типа без паники:
s, ok := i.(string)
fmt.Println(s, ok)  // hello true

n, ok := i.(int)
fmt.Println(n, ok)  // 0 false (нулевое значение int, ok = false)
```

`ok == false` — безопасно, нет паники, `value` = нулевое значение типа.

---

## Практические примеры

### Извлечение конкретного типа для доп. функциональности

```go
type Animal interface {
    Speak() string
}

type Dog struct{ Name string }
func (d Dog) Speak() string { return "Гав!" }
func (d Dog) Fetch() string { return d.Name + " принёс мяч!" }

type Cat struct{ Name string }
func (c Cat) Speak() string { return "Мяу!" }
func (c Cat) Purr() string  { return c.Name + " мурлычет..." }

func interact(a Animal) {
    fmt.Println(a.Speak())

    // Попытка использовать Dog-специфичный метод:
    if dog, ok := a.(Dog); ok {
        fmt.Println(dog.Fetch())
    }

    // Попытка использовать Cat-специфичный метод:
    if cat, ok := a.(Cat); ok {
        fmt.Println(cat.Purr())
    }
}

interact(Dog{Name: "Рекс"})
// Гав!
// Рекс принёс мяч!

interact(Cat{Name: "Мурка"})
// Мяу!
// Мурка мурлычет...
```

### Работа с JSON (dynamic values)

```go
import "encoding/json"

// json.Unmarshal в interface{} — все числа станут float64
var data interface{}
json.Unmarshal([]byte(`{"name":"Alice","age":30}`), &data)

m, ok := data.(map[string]interface{})
if !ok {
    return
}

name, _ := m["name"].(string)
age, _ := m["age"].(float64)  // числа в JSON → float64!
fmt.Println(name, int(age))    // Alice 30
```

### Проверка реализации опциональных интерфейсов

```go
type Flusher interface {
    Flush() error
}

func writeAndFlush(w io.Writer, data []byte) error {
    _, err := w.Write(data)
    if err != nil {
        return err
    }

    // Если Writer поддерживает Flush — вызываем:
    if flusher, ok := w.(Flusher); ok {
        return flusher.Flush()
    }

    return nil
}
```

Это паттерн проверки «опциональных возможностей» — часто встречается в стандартной библиотеке.

---

## Цепочка type assertions

```go
type Stringer interface{ String() string }
type Closer  interface{ Close() error }

func processResource(res interface{}) {
    if sc, ok := res.(interface {
        String() string
        Close() error
    }); ok {
        fmt.Println(sc.String())
        defer sc.Close()
    }
}
```

Можно делать assertion к анонимному интерфейсу прямо в коде.

---

## Типичные ошибки

**Ошибка 1**: Использовать простую форму без уверенности в типе.

```go
func process(v interface{}) {
    s := v.(string)  // паника если v не string!
    fmt.Println(s)
}

process("hello")  // OK
process(42)       // panic!
```

**Ошибка 2**: Пытаться сделать assertion к несовместимому типу.

```go
var r io.Reader = os.Stdin

// Можно — os.Stdin это *os.File, который реализует io.Writer:
w, ok := r.(io.Writer)
fmt.Println(ok)  // true

// Нельзя сделать assertion к конкретному типу другого пакета напрямую:
// ... но к интерфейсу — можно, Go проверит во runtime
```

**Ошибка 3**: Assertion к nil интерфейсу.

```go
var i interface{} = nil
s, ok := i.(string)
fmt.Println(s, ok)  // "" false — это нормально, не паника

// НО простая форма к nil — паника:
s2 := i.(string)  // panic: interface conversion: interface is nil, not string
```

---

## Итог

- `v.(T)` — извлечь значение типа T; паника если тип неверный
- `v, ok := i.(T)` — безопасная форма; при неверном типе ok=false, паники нет
- Используй comma-ok форму в большинстве случаев
- Простую форму — только когда тип гарантирован (например, сразу после type switch)
- Type assertion к интерфейсу проверяет, реализует ли конкретный тип этот интерфейс
