---
title: "Встраивание структур"
description: "Встраивание, продвижение методов, отличие от наследования, конфликт имён"
order: 3
---

# Встраивание структур

Встраивание (embedding) — механизм Go для повторного использования кода. Это **не наследование**, хотя поверхностно похоже. Понимание разницы критично для правильного проектирования.

## Продвижение полей и методов

Когда структура встраивает другую, все поля и методы встроенной «продвигаются» во внешнюю:

```go
type Animal struct {
    Name string
    Age  int
}

func (a Animal) Describe() string {
    return fmt.Sprintf("%s, %d лет", a.Name, a.Age)
}

type Dog struct {
    Animal        // встраивание
    Breed  string
}

func main() {
    d := Dog{
        Animal: Animal{Name: "Рекс", Age: 3},
        Breed:  "Лабрадор",
    }

    // Доступ через имя встроенного типа:
    fmt.Println(d.Animal.Name)      // Рекс
    fmt.Println(d.Animal.Describe()) // Рекс, 3 лет

    // Продвижение — через Dog напрямую:
    fmt.Println(d.Name)      // Рекс (d.Animal.Name)
    fmt.Println(d.Describe()) // Рекс, 3 лет (d.Animal.Describe())
    fmt.Println(d.Breed)     // Лабрадор
}
```

---

## Отличие от наследования

В ООП-языках наследование создаёт отношение «является» (is-a). В Go встраивание — это «содержит» (has-a) с возможностью делегирования:

```go
type Writer struct{}
func (w Writer) Write(data []byte) (int, error) { ... }

type BufferedWriter struct {
    Writer      // встраивание, не наследование
    buf []byte
}

// Dog НЕ является Animal в смысле полиморфизма:
var a Animal = d    // ОШИБКА: Dog не реализует Animal (нет такого интерфейса)

// Но поля и методы доступны через продвижение
```

**Полиморфизм в Go** достигается через интерфейсы, не через встраивание:

```go
type Describer interface {
    Describe() string
}

// Dog имеет метод Describe() через продвижение — значит, удовлетворяет интерфейсу
var d Describer = Dog{Animal: Animal{Name: "Рекс"}}
fmt.Println(d.Describe())  // работает!
```

> 💡 Интерфейсы и их мощь подробно разберём в главе 7.

---

## Переопределение (shadowing) методов

Внешняя структура может «перекрыть» продвинутый метод:

```go
type Base struct{}

func (b Base) Method() string { return "Base.Method" }

type Derived struct {
    Base
}

func (d Derived) Method() string { return "Derived.Method" }

d := Derived{}
fmt.Println(d.Method())        // Derived.Method — собственный метод
fmt.Println(d.Base.Method())   // Base.Method — явный доступ к встроенному
```

---

## Конфликт имён

Если несколько встроенных структур имеют поля/методы с одинаковым именем — нужно обращаться явно:

```go
type Logger struct {
    Level string
}
func (l Logger) Log(msg string) { fmt.Println("[" + l.Level + "]", msg) }

type Database struct {
    Level string  // конфликт имён с Logger.Level!
}
func (db Database) Log(msg string) { fmt.Println("[DB]", msg) }

type Service struct {
    Logger
    Database
}

s := Service{
    Logger:   Logger{Level: "INFO"},
    Database: Database{Level: "DEBUG"},
}

// s.Level  — ОШИБКА: ambiguous selector s.Level
// s.Log()  — ОШИБКА: ambiguous selector s.Log

// Нужно явное обращение:
fmt.Println(s.Logger.Level)   // INFO
fmt.Println(s.Database.Level) // DEBUG
s.Logger.Log("запуск")        // [INFO] запуск
s.Database.Log("запрос")      // [DB] запрос
```

---

## Встраивание указателя

Можно встраивать не только значение, но и указатель:

```go
type Engine struct {
    Power int
}
func (e *Engine) Start() { fmt.Println("Двигатель запущен") }

type Car struct {
    *Engine  // встраиваем указатель
    Model string
}

car := Car{
    Engine: &Engine{Power: 200},
    Model:  "Lada",
}

car.Start()  // продвижение работает через указатель
fmt.Println(car.Power)  // 200

// Осторожно: если Engine == nil, вызов Start() вызовет панику
var brokenCar Car
// brokenCar.Start()  // паника! *Engine == nil
```

---

## Практический пример: составной логгер

```go
package main

import (
    "fmt"
    "time"
)

type TimestampMixin struct{}

func (t TimestampMixin) Now() string {
    return time.Now().Format("15:04:05")
}

type PrefixMixin struct {
    Prefix string
}

func (p PrefixMixin) FormatMsg(msg string) string {
    return fmt.Sprintf("[%s] %s", p.Prefix, msg)
}

type Logger struct {
    TimestampMixin
    PrefixMixin
}

func (l Logger) Log(msg string) {
    fmt.Printf("%s %s\n", l.Now(), l.FormatMsg(msg))
}

func main() {
    log := Logger{
        PrefixMixin: PrefixMixin{Prefix: "APP"},
    }

    log.Log("сервер запущен")
    // 14:32:01 [APP] сервер запущен
}
```

---

## Встраивание интерфейсов в структуры

Интерфейс тоже можно встроить в структуру — полезно для mock-объектов:

```go
type Storer interface {
    Save(data []byte) error
    Load(id string) ([]byte, error)
}

type LoggingStorer struct {
    Storer  // встроенный интерфейс
}

func (ls *LoggingStorer) Save(data []byte) error {
    fmt.Println("сохраняем", len(data), "байт")
    return ls.Storer.Save(data)  // делегируем встроенному
}
// Load() делегируется автоматически
```

---

## Итог

- Встраивание = содержит (has-a), не наследование (is-a)
- Поля и методы встроенной структуры «продвигаются» во внешнюю
- Внешняя структура может перекрыть продвинутый метод
- При конфликте имён — нужно явное обращение: `s.Embedded.Field`
- Встраивание реализует делегирование, а полиморфизм достигается через интерфейсы
- Можно встраивать как значение, так и указатель (`*T`)
