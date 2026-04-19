---
title: "Псевдонимы и определения типов"
description: "type MyType = int vs type MyType int — принципиальная разница и когда использовать каждый"
order: 5
---

# Псевдонимы и определения типов

В Go есть два способа создать новое имя для типа, и они принципиально отличаются. Путаница между ними — источник реальных ошибок.

## Два разных синтаксиса

```go
type Celsius float64      // определение нового типа (type definition)
type MyFloat = float64    // псевдоним (type alias)
```

Выглядят похоже, но работают по-разному.

---

## Type Definition — новый тип

```go
type Celsius float64
type Fahrenheit float64
```

Создаёт **новый, отдельный тип** на основе существующего. Новый тип:
- Имеет тот же набор операций, что и базовый (арифметика для числовых типов)
- **Не взаимозаменяем** с базовым типом без явного преобразования
- Может иметь свои методы

```go
type Celsius float64
type Fahrenheit float64

func CToF(c Celsius) Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

func FToC(f Fahrenheit) Celsius {
    return Celsius((f - 32) * 5 / 9)
}

func main() {
    boiling := Celsius(100)
    fmt.Printf("%.1f°C = %.1f°F\n", boiling, CToF(boiling))
    // 100.0°C = 212.0°F

    // Нельзя смешать Celsius и Fahrenheit без явного преобразования:
    var temp float64 = 100
    c := Celsius(temp)         // OK: явное преобразование
    // c2 := Celsius(boiling)  // OK: одинаковый тип
    // c3 := temp              // ОШИБКА: cannot use temp (float64) as Celsius
}
```

### Методы на новом типе

Это главная причина использовать type definition — можно добавить методы:

```go
type Celsius float64

func (c Celsius) String() string {
    return fmt.Sprintf("%.2f°C", float64(c))
}

func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

func main() {
    temp := Celsius(37.5)
    fmt.Println(temp)              // 37.50°C (метод String() вызывается автоматически)
    fmt.Println(temp.ToFahrenheit())  // 99.50°F
}
```

> 💡 Методы на типах подробно разберём в главе 6 «Указатели и структуры».

### Семантическая безопасность

Type definition помогает предотвратить логические ошибки:

```go
type UserID int64
type ProductID int64

func GetUser(id UserID) (*User, error) { ... }
func GetProduct(id ProductID) (*Product, error) { ... }

func main() {
    userID := UserID(42)
    productID := ProductID(42)

    GetUser(productID)    // ОШИБКА компиляции: нельзя передать ProductID туда, где ожидается UserID
    GetUser(userID)       // OK
}
```

Без отдельных типов `GetUser(productID)` скомпилировался бы без ошибок — и мы бы получали данные не того объекта.

---

## Type Alias — псевдоним

```go
type MyFloat = float64
```

Создаёт **другое имя для того же типа**. Это буквально синоним — новый тип не создаётся. Псевдоним:
- Полностью взаимозаменяем с оригинальным типом
- **Не может иметь новых методов**
- Используется для совместимости и рефакторинга

```go
type MyFloat = float64

var x MyFloat = 3.14
var y float64 = x  // OK! MyFloat и float64 — один тип
fmt.Println(x + y) // 6.28

// Попытка добавить метод — ошибка:
// func (f MyFloat) String() string { return ... }
// ОШИБКА: cannot define new methods on non-local type float64
```

### Главное применение: рефакторинг и совместимость

Предположим, ты переносишь тип из одного пакета в другой:

```go
// Старый код в пакете old:
package old

type Config struct { ... }

// Новый код в пакете new:
package new

type Config struct { ... }

// В пакете old добавляем псевдоним для сохранения совместимости:
package old

import "myapp/new"

type Config = new.Config  // псевдоним: старое имя ссылается на новый тип
```

Теперь весь код, который использовал `old.Config`, продолжает работать без изменений, а новый код может использовать `new.Config`.

### Стандартная библиотека: byte и rune

```go
// В spec Go объявлены как псевдонимы:
type byte = uint8
type rune = int32
```

Это настоящие псевдонимы — `byte` и `uint8` полностью взаимозаменяемы, `rune` и `int32` тоже:

```go
var b byte = 65
var u uint8 = b  // OK, один тип

var r rune = 'А'
var i int32 = r  // OK, один тип
```

---

## Сравнительная таблица

| Свойство | `type T int` | `type T = int` |
|----------|-------------|----------------|
| Новый тип создаётся? | Да | Нет (синоним) |
| Взаимозаменяем с оригиналом? | Нет (нужно явное приведение) | Да |
| Можно добавить методы? | Да | Нет |
| Семантическая безопасность | Есть | Нет |
| Основное применение | Доменные типы, enum | Рефакторинг, совместимость |

---

## Практические примеры

### Типы для идентификаторов (type definition)

```go
type UserID int64
type OrderID int64
type ProductID int64

type Order struct {
    ID         OrderID
    UserID     UserID
    ProductIDs []ProductID
}

// Теперь нельзя случайно передать OrderID туда, где нужен UserID
```

### Enum-подобный тип

```go
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
    StatusDeleted
)

func (s Status) String() string {
    switch s {
    case StatusPending:
        return "pending"
    case StatusActive:
        return "active"
    case StatusInactive:
        return "inactive"
    case StatusDeleted:
        return "deleted"
    default:
        return fmt.Sprintf("unknown(%d)", int(s))
    }
}

func processUser(status Status) {
    if status == StatusActive {
        fmt.Println("Обрабатываем активного пользователя")
    }
}
```

### Единицы измерения

```go
type Meters float64
type Seconds float64
type MetersPerSecond float64

func Speed(distance Meters, time Seconds) MetersPerSecond {
    return MetersPerSecond(float64(distance) / float64(time))
}

func main() {
    distance := Meters(100)
    time := Seconds(9.58)
    speed := Speed(distance, time)
    fmt.Printf("Скорость: %.2f м/с\n", float64(speed))
    // Скорость: 10.44 м/с
}
```

---

## Типичные ошибки

**Ошибка 1**: Ожидать, что type definition совместим с базовым типом.

```go
type MyInt int

func add(a, b int) int { return a + b }

x := MyInt(5)
y := MyInt(3)

// add(x, y)  // ОШИБКА: cannot use x (type MyInt) as type int
add(int(x), int(y))  // OK: явное преобразование
```

**Ошибка 2**: Попытка добавить методы к псевдониму.

```go
type MyString = string

// func (s MyString) Upper() string { return strings.ToUpper(string(s)) }
// ОШИБКА: cannot define new methods on non-local type string

// Правильно — использовать type definition:
type MyString string

func (s MyString) Upper() string { return strings.ToUpper(string(s)) }
```

---

## Итог

- `type T int` — **новый тип**: отдельный, несовместимый с `int`, может иметь методы
- `type T = int` — **псевдоним**: синоним для `int`, полностью взаимозаменяем
- Type definition используй для семантической безопасности (UserID vs ProductID) и доменных типов
- Type alias используй для рефакторинга и обратной совместимости
- `byte = uint8` и `rune = int32` — встроенные псевдонимы в стандартной библиотеке
