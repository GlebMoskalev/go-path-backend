---
title: "Методы"
description: "Value receiver vs pointer receiver, правила выбора, method set, методы на не-struct типах"
order: 4
---

# Методы

Метод — функция с получателем (receiver). В Go нет классов, но любой именованный тип может иметь методы.

## Синтаксис

```go
func (получатель ТипПолучателя) ИмяМетода(параметры) возвращаемыеТипы {
    // тело
}
```

```go
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

rect := Rectangle{Width: 10, Height: 5}
fmt.Println(rect.Area())      // 50
fmt.Println(rect.Perimeter()) // 30
```

---

## Value Receiver vs Pointer Receiver

Это одно из самых важных решений при написании методов.

### Value Receiver — `(r Rectangle)`

Метод получает **копию** значения. Изменения не видны снаружи:

```go
func (r Rectangle) ScaleWrong(factor float64) {
    r.Width *= factor   // изменяем копию
    r.Height *= factor
}

rect := Rectangle{10, 5}
rect.ScaleWrong(2)
fmt.Println(rect)  // {10 5} — не изменился!
```

### Pointer Receiver — `(r *Rectangle)`

Метод получает **указатель** на оригинал. Изменения видны снаружи:

```go
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor   // изменяем оригинал
    r.Height *= factor
}

rect := Rectangle{10, 5}
rect.Scale(2)          // Go автоматически берёт &rect
fmt.Println(rect)      // {20 10} — изменился!
```

Go автоматически преобразует `rect.Scale(2)` в `(&rect).Scale(2)` — не нужно явно писать `&`.

---

## Правила выбора receiver

**Используй pointer receiver, если:**
1. Метод изменяет состояние объекта
2. Структура большая (копирование дорого)
3. Нужна согласованность (если есть хоть один pointer receiver, лучше все сделать pointer)

**Используй value receiver, если:**
1. Метод только читает данные
2. Тип — маленький, дешёвый для копирования (int, float64, small struct)
3. Тип предназначен для неизменяемого использования

```go
type Counter struct {
    n int
}

// Pointer receiver — изменяет состояние:
func (c *Counter) Increment() { c.n++ }
func (c *Counter) Reset()     { c.n = 0 }

// Value receiver — только читает:
// НО: для согласованности лучше тоже pointer receiver
func (c *Counter) Value() int { return c.n }

// Маленький тип — value receiver оправдан:
type Point struct{ X, Y float64 }
func (p Point) Distance() float64 { return math.Sqrt(p.X*p.X + p.Y*p.Y) }
func (p Point) String() string    { return fmt.Sprintf("(%.2f, %.2f)", p.X, p.Y) }
```

---

## Method Set

Method set типа определяет, какие интерфейсы он реализует.

| Тип | Доступные методы |
|-----|-----------------|
| `T` | методы с value receiver |
| `*T` | методы с value receiver + pointer receiver |

```go
type Greeter interface {
    Greet() string
}

type Person struct {
    Name string
}

func (p Person) Greet() string {
    return "Привет, я " + p.Name
}

// Person реализует Greeter через value receiver:
var g Greeter = Person{Name: "Алиса"}  // OK
var g2 Greeter = &Person{Name: "Боб"}  // тоже OK — *Person включает value методы
```

Но если метод с pointer receiver:

```go
type Writer interface {
    Write(data string)
}

func (p *Person) Write(data string) {
    fmt.Println(p.Name, "пишет:", data)
}

var w Writer = &Person{Name: "Алиса"}  // OK
// var w2 Writer = Person{Name: "Боб"}  // ОШИБКА: Person не реализует Writer
// (метод Write есть только у *Person)
```

---

## Методы на не-struct типах

Методы можно определять для любого именованного типа в том же пакете:

```go
// Методы на числовом типе:
type Celsius float64
type Fahrenheit float64

func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}

func (c Celsius) String() string {
    return fmt.Sprintf("%.2f°C", float64(c))
}

temp := Celsius(100)
fmt.Println(temp)                // 100.00°C (вызывается String())
fmt.Println(temp.ToFahrenheit()) // 212°F

// Методы на слайсе:
type StringSlice []string

func (s StringSlice) Contains(target string) bool {
    for _, v := range s {
        if v == target {
            return true
        }
    }
    return false
}

func (s StringSlice) Filter(pred func(string) bool) StringSlice {
    result := make(StringSlice, 0)
    for _, v := range s {
        if pred(v) {
            result = append(result, v)
        }
    }
    return result
}

words := StringSlice{"apple", "banana", "cherry", "avocado"}
fmt.Println(words.Contains("banana"))  // true
aWords := words.Filter(func(s string) bool { return s[0] == 'a' })
fmt.Println(aWords)  // [apple avocado]
```

Нельзя добавлять методы к чужим типам из другого пакета:

```go
// Нельзя:
func (s string) Upper() string { ... }  // ОШИБКА: cannot define methods on non-local type

// Правильно — создай тип-обёртку:
type MyString string
func (s MyString) Upper() MyString { return MyString(strings.ToUpper(string(s))) }
```

---

## Цепочка вызовов (method chaining)

Возвращая receiver из методов, можно строить цепочки:

```go
type QueryBuilder struct {
    table  string
    wheres []string
    limit  int
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
    q.table = table
    return q
}

func (q *QueryBuilder) Where(condition string) *QueryBuilder {
    q.wheres = append(q.wheres, condition)
    return q
}

func (q *QueryBuilder) Limit(n int) *QueryBuilder {
    q.limit = n
    return q
}

func (q *QueryBuilder) Build() string {
    query := "SELECT * FROM " + q.table
    if len(q.wheres) > 0 {
        query += " WHERE " + strings.Join(q.wheres, " AND ")
    }
    if q.limit > 0 {
        query += fmt.Sprintf(" LIMIT %d", q.limit)
    }
    return query
}

query := (&QueryBuilder{}).
    From("users").
    Where("active = true").
    Where("age > 18").
    Limit(10).
    Build()

fmt.Println(query)
// SELECT * FROM users WHERE active = true AND age > 18 LIMIT 10
```

---

## Итог

- Метод = функция с receiver: `func (r T) Method() {}`
- Value receiver — получает копию, изменения не видны снаружи
- Pointer receiver — получает указатель, изменения видны снаружи
- Используй pointer receiver если: изменяешь состояние, большая структура, нужна согласованность
- `T` имеет value методы; `*T` имеет и value, и pointer методы
- Методы можно определять для любого именованного типа в своём пакете, не только struct
