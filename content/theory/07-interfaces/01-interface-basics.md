---
title: "Основы интерфейсов"
description: "Неявная реализация, duck typing, интерфейс как контракт, пустой интерфейс any"
order: 1
---

# Основы интерфейсов

Интерфейсы — самая мощная концепция Go. Они обеспечивают полиморфизм без иерархий наследования, через чистый duck typing.

## Объявление интерфейса

```go
type Имя interface {
    Метод1(параметры) возвращаемыеТипы
    Метод2(параметры) возвращаемыеТипы
}
```

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

---

## Неявная реализация

В Go нет ключевого слова `implements`. Тип реализует интерфейс **автоматически**, если имеет все нужные методы:

```go
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Оба типа реализуют Shape — без явного объявления:
var s1 Shape = Rectangle{10, 5}
var s2 Shape = Circle{7}

func printShape(s Shape) {
    fmt.Printf("Площадь: %.2f, Периметр: %.2f\n", s.Area(), s.Perimeter())
}

printShape(s1)  // Площадь: 50.00, Периметр: 30.00
printShape(s2)  // Площадь: 153.94, Периметр: 43.98
```

---

## Duck Typing

«Если это ходит как утка и крякает как утка — это утка». Тип удовлетворяет интерфейсу, если имеет нужные методы, независимо от того, знает ли он об этом интерфейсе:

```go
// Стандартный интерфейс из пакета fmt:
type Stringer interface {
    String() string
}

// Наш тип — ничего не знает о fmt.Stringer
type Temperature struct {
    Celsius float64
}

func (t Temperature) String() string {
    return fmt.Sprintf("%.1f°C", t.Celsius)
}

// Но т.к. у него есть метод String(), fmt.Println его вызовет:
temp := Temperature{36.6}
fmt.Println(temp)  // 36.6°C
```

Это позволяет типам из разных пакетов работать с интерфейсами, которые они не видели при написании.

---

## Интерфейс как контракт

Интерфейс определяет **что** должен делать тип, не **как**. Это позволяет подменять реализации:

```go
type Notifier interface {
    Send(to, message string) error
}

// Реализация через email:
type EmailNotifier struct {
    SMTPHost string
}
func (e EmailNotifier) Send(to, msg string) error {
    fmt.Printf("Email to %s: %s\n", to, msg)
    return nil
}

// Реализация через SMS:
type SMSNotifier struct {
    APIKey string
}
func (s SMSNotifier) Send(to, msg string) error {
    fmt.Printf("SMS to %s: %s\n", to, msg)
    return nil
}

// Реализация для тестов:
type MockNotifier struct {
    Sent []string
}
func (m *MockNotifier) Send(to, msg string) error {
    m.Sent = append(m.Sent, to+": "+msg)
    return nil
}

// Функция работает с любой реализацией:
func notifyUser(n Notifier, user, msg string) error {
    return n.Send(user, msg)
}

notifyUser(EmailNotifier{}, "alice@example.com", "Привет!")
notifyUser(SMSNotifier{}, "+79991234567", "Код: 1234")
```

---

## Пустой интерфейс: any

`any` (псевдоним `interface{}`) — интерфейс без методов. Реализуется **любым** типом:

```go
var x any = 42
x = "hello"
x = []int{1, 2, 3}
x = true
```

Используется когда тип неизвестен заранее:

```go
func printAny(v any) {
    fmt.Printf("Тип: %T, Значение: %v\n", v, v)
}

printAny(42)          // Тип: int, Значение: 42
printAny("hello")     // Тип: string, Значение: hello
printAny([]int{1,2})  // Тип: []int, Значение: [1 2]
```

### Когда использовать any

**Уместно:**
- Контейнеры с разнородными типами
- `fmt.Println(args ...any)` — функция с разными типами аргументов
- Работа с JSON при неизвестной структуре

**Избегай:**
- Не используй `any` вместо generics когда типы известны
- Не используй `any` чтобы обойти систему типов — это потеря type safety

```go
// Плохо: using any when generics better
func maxAny(a, b any) any {
    // нужны runtime-проверки типов...
}

// Хорошо: generics (Go 1.18+)
func max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

---

## Составные интерфейсы

Интерфейсы можно объединять:

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

// Тип реализует ReadWriter, если реализует и Read, и Write
```

```go
type Closer interface {
    Close() error
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

---

## Принцип маленьких интерфейсов

Go-сообщество придерживается правила: **интерфейсы должны быть маленькими**. Лучший интерфейс — из одного метода:

```go
// Хорошо — маленький, точный контракт:
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Менее хорошо — большой интерфейс с множеством методов:
type FileManager interface {
    Open(path string) error
    Close() error
    Read(p []byte) (int, error)
    Write(p []byte) (int, error)
    Seek(offset int64, whence int) (int64, error)
    Stat() (FileInfo, error)
    // ...
}
```

Маленький интерфейс проще реализовать и замокировать в тестах.

---

## Проверка реализации интерфейса

Иногда нужно убедиться в compile-time, что тип реализует интерфейс:

```go
// Compile-time проверка:
var _ Shape = (*Rectangle)(nil)  // не выделяет память, просто проверка типов
var _ Shape = Circle{}

// Если Rectangle не реализует Shape — ошибка компиляции
```

---

## Итог

- Интерфейс — набор методов, описывающий поведение
- Реализация неявная: достаточно иметь нужные методы
- Duck typing: «если есть нужные методы — тип подходит»
- `any` (≡ `interface{}`) — пустой интерфейс, принимает всё
- Предпочитай маленькие интерфейсы из 1-3 методов
- Составные интерфейсы (`Reader` + `Writer` = `ReadWriter`) строятся встраиванием
