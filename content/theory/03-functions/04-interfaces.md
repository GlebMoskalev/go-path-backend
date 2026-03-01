---
title: "Интерфейсы"
description: "Неявная реализация интерфейсов, пустой интерфейс, type assertion, type switch"
order: 4
---

# Интерфейсы

Интерфейс — набор сигнатур методов. Тип реализует интерфейс **неявно**, просто имея нужные методы.

## Объявление

```go
type Speaker interface {
    Speak() string
}
```

## Неявная реализация

```go
type Dog struct{ Name string }

func (d Dog) Speak() string {
    return d.Name + ": Гав!"
}

type Cat struct{ Name string }

func (c Cat) Speak() string {
    return c.Name + ": Мяу!"
}
```

Оба типа реализуют `Speaker` без явного указания:

```go
func greet(s Speaker) {
    fmt.Println(s.Speak())
}

greet(Dog{Name: "Бобик"})  // Бобик: Гав!
greet(Cat{Name: "Мурка"})  // Мурка: Мяу!
```

## Полиморфизм

```go
animals := []Speaker{
    Dog{Name: "Рекс"},
    Cat{Name: "Барсик"},
    Dog{Name: "Шарик"},
}

for _, a := range animals {
    fmt.Println(a.Speak())
}
```

## Пустой интерфейс

`interface{}` (или `any` в Go 1.18+) принимает значение любого типа:

```go
func printAnything(v any) {
    fmt.Printf("тип: %T, значение: %v\n", v, v)
}

printAnything(42)
printAnything("hello")
printAnything([]int{1, 2, 3})
```

## Type assertion

Извлечение конкретного типа из интерфейса:

```go
var s Speaker = Dog{Name: "Рекс"}

// Безопасная форма
dog, ok := s.(Dog)
if ok {
    fmt.Println("Это собака:", dog.Name)
}

// Опасная форма — паника если тип не совпадает
dog := s.(Dog)
```

## Type switch

```go
func describe(v any) string {
    switch val := v.(type) {
    case int:
        return fmt.Sprintf("целое число: %d", val)
    case string:
        return fmt.Sprintf("строка: %q", val)
    case bool:
        if val {
            return "истина"
        }
        return "ложь"
    default:
        return fmt.Sprintf("неизвестный тип: %T", val)
    }
}
```

## Композиция интерфейсов

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

## Популярные интерфейсы стандартной библиотеки

| Интерфейс      | Пакет    | Методы                               |
|----------------|----------|--------------------------------------|
| `io.Reader`    | `io`     | `Read([]byte) (int, error)`          |
| `io.Writer`    | `io`     | `Write([]byte) (int, error)`         |
| `fmt.Stringer` | `fmt`    | `String() string`                    |
| `error`        | builtin  | `Error() string`                     |
| `sort.Interface` | `sort` | `Len()`, `Less()`, `Swap()`         |

## Stringer

```go
type User struct {
    Name string
    Age  int
}

func (u User) String() string {
    return fmt.Sprintf("%s (%d лет)", u.Name, u.Age)
}

u := User{Name: "Гофер", Age: 15}
fmt.Println(u) // Гофер (15 лет)
```

## Nil-интерфейс

```go
var s Speaker // nil — нет значения и нет типа

// Интерфейс с nil-значением
var d *Dog = nil
var s Speaker = d // s != nil! (тип есть, значение nil)
```

> **Совет:** проектируйте маленькие интерфейсы (1-3 метода). Большие интерфейсы сложнее реализовать и тестировать.
