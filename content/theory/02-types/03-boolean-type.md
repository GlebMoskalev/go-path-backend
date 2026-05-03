---
title: "Булевый тип"
description: "bool, ленивые вычисления (&&, ||), идиоматические паттерны использования"
order: 3
---

# Булевый тип

Тип `bool` — один из самых простых в Go, но ленивые вычисления операторов `&&` и `||` (short-circuit evaluation) и идиоматические паттерны его использования стоит разобрать отдельно.

## Основы

```go
var b bool         // нулевое значение: false
b = true
b = false

active := true
deleted := false

fmt.Printf("Тип: %T, значение: %v\n", active, active)
// Тип: bool, значение: true
```

В Go нет неявного преобразования в bool. Это принципиальное решение:

```go
n := 0
if n {           // ОШИБКА: non-bool n (type int) used as if condition
}
if n != 0 {      // OK — явное сравнение
}

s := ""
if s {           // ОШИБКА
}
if s != "" {     // OK
}
if len(s) > 0 { // OK — явная длина
}
```

---

## Логические операторы

```go
a, b := true, false

fmt.Println(a && b)   // false — AND
fmt.Println(a || b)   // true  — OR
fmt.Println(!a)       // false — NOT
```

### Ленивые вычисления (short-circuit evaluation)

`&&` и `||` не всегда вычисляют оба операнда — если результат уже известен по левому, правый пропускается. Это не оптимизация, а гарантия языка.

**`&&` (AND)**: если левый операнд `false` — правый не вычисляется.

```go
func riskyOp() bool {
    fmt.Println("riskyOp вызван")
    return true
}

x := 0
if x != 0 && riskyOp() {  // riskyOp НЕ вызовется
    fmt.Println("условие истинно")
}
// Вывод: ничего — riskyOp не был вызван
```

**`||` (OR)**: если левый операнд `true` — правый не вычисляется.

```go
if x == 0 || riskyOp() {  // riskyOp НЕ вызовется
    fmt.Println("условие истинно")
}
// Вывод: условие истинно — riskyOp не был вызван
```

### Практическое применение short-circuit

**Защита от nil pointer:**

> `nil` и указатели подробно разберём позже — в разделе про указатели. Пока достаточно знать: `nil` означает «нет значения», и обращение к методу через `nil` вызывает панику.

```go
// Безопасно: если user == nil, метод не вызовется
if user != nil && user.IsActive() {
    sendEmail(user)
}
```

**Проверка длины перед индексацией:**

```go
items := []string{"a", "b"}
// Безопасно: если len(items) == 0, items[0] не вычисляется
if len(items) > 0 && items[0] == "a" {
    fmt.Println("первый элемент — a")
}
```

**Проверка ошибки перед использованием результата:**

```go
file, err := os.Open("data.txt")
// Если err != nil, операция с file не выполняется
if err == nil && fileSize(file) > 0 {
    process(file)
}
```

---

## Типичные паттерны

### Функция, возвращающая bool

```go
func isEmpty(s string) bool {
    return s == ""
}

func isValidAge(age int) bool {
    return age >= 0 && age <= 150
}

func isLeapYear(year int) bool {
    return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}
```

### Comma-ok идиома

Многие операции в Go возвращают пару `(значение, bool)`:

```go
m := map[string]int{"a": 1, "b": 2}

// Без проверки — получим нулевое значение для отсутствующего ключа:
val := m["c"]  // 0 — но c нет в карте!

// С проверкой — comma-ok:
val, ok := m["c"]
if ok {
    fmt.Println("найдено:", val)
} else {
    fmt.Println("ключ не найден")
}
```

> 💡 Comma-ok широко используется для map, type assertion и channel-операций. Подробно в главах 5 (коллекции) и 7 (интерфейсы).

### Булевы флаги в структурах

```go
type User struct {
    Name    string
    Email   string
    Active  bool
    Admin   bool
    Deleted bool
}

func (u *User) CanLogin() bool {
    return u.Active && !u.Deleted
}

func (u *User) CanManageUsers() bool {
    return u.CanLogin() && u.Admin
}
```

### Построение сложных условий

```go
func isValidOrder(order Order) bool {
    hasItems := len(order.Items) > 0
    validAmount := order.Total > 0 && order.Total < 1_000_000
    hasDelivery := order.Address != "" || order.PickupPoint != ""

    return hasItems && validAmount && hasDelivery
}
```

Разбивка на именованные переменные делает логику читаемой — лучше, чем одно длинное условие.

---

## Сравнение с другими языками

В Python, JavaScript, Ruby любое значение может быть truthy/falsy. В Go — нет. Это намеренное ограничение:

```go
// Python: if s: — работает для непустой строки
// Go: так нельзя — нужно явное условие

s := "hello"
if s != "" { // явно
    fmt.Println("непустая строка")
}

n := 42
if n > 0 { // явно
    fmt.Println("положительное число")
}
```

Это делает код чуть более многословным, но намерение всегда явно.

---

## Итог

- `bool` принимает только `true` или `false`, нулевое значение — `false`
- Нет неявного приведения к bool — все условия должны быть явными
- Short-circuit evaluation: `&&` — пропускает правый операнд при `false` слева, `||` — при `true` слева
- Это позволяет безопасно писать `ptr != nil && ptr.Method()`
- Comma-ok идиома (`val, ok := map[key]`) — стандартный способ обработки опциональных значений
