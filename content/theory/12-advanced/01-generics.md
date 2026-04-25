---
title: "Дженерики"
description: "Типовые параметры, ограничения (constraints), когда использовать"
order: 1
---

# Дженерики

Дженерики (обобщённые типы) появились в Go 1.18. Они позволяют писать функции и типы, работающие с несколькими типами данных без дублирования кода.

## Проблема до дженериков

```go
// Без дженериков: отдельная функция для каждого типа
func MaxInt(a, b int) int {
    if a > b { return a }
    return b
}

func MaxFloat64(a, b float64) float64 {
    if a > b { return a }
    return b
}

// Или: потеря типобезопасности через any
func Max(a, b any) any {
    // нет возможности сравнить без type assertion
}
```

---

## Синтаксис типовых параметров

```go
// T — типовой параметр, comparable — ограничение
func Contains[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

// Использование:
Contains([]int{1, 2, 3}, 2)          // true
Contains([]string{"a", "b"}, "c")    // false
// тип T выводится автоматически

// Явное указание типа:
Contains[int]([]int{1, 2, 3}, 2)
```

---

## Ограничения (constraints)

Ограничение определяет какие типы допустимы для типового параметра.

### Встроенные ограничения

```go
// any = interface{} — любой тип
func Print[T any](v T) { fmt.Println(v) }

// comparable — типы поддерживающие == и !=
func Equal[T comparable](a, b T) bool { return a == b }
```

### Пакет golang.org/x/exp/constraints (стандартизируется)

```go
import "golang.org/x/exp/constraints"

// constraints.Ordered: int, float64, string и их вариации
func Max[T constraints.Ordered](a, b T) T {
    if a > b { return a }
    return b
}

Max(3, 5)       // int → 5
Max(3.14, 2.7)  // float64 → 3.14
Max("b", "a")   // string → "b"
```

### Встроенные в stdlib (Go 1.21+)

```go
import "cmp"

// cmp.Ordered — int, float, string
func Min[T cmp.Ordered](a, b T) T {
    if a < b { return a }
    return b
}
```

### Пользовательские ограничения

```go
// Ограничение через union type
type Number interface {
    int | int8 | int16 | int32 | int64 |
    float32 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

Sum([]int{1, 2, 3})       // 6
Sum([]float64{1.1, 2.2})  // 3.3
```

Тильда `~` — включить типы с тем же underlying type:

```go
type Celsius float64
type Fahrenheit float64

type Temperature interface {
    ~float32 | ~float64  // включает Celsius, Fahrenheit и другие float-based типы
}

func AbsDiff[T Temperature](a, b T) T {
    if a > b { return a - b }
    return b - a
}

AbsDiff(Celsius(100), Celsius(37))  // работает
```

---

## Обобщённые типы

```go
// Обобщённый Stack
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    last := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return last, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// Использование:
s := Stack[int]{}
s.Push(1)
s.Push(2)
v, ok := s.Pop()  // v=2, ok=true
```

---

## Утилитарные функции

```go
// Map: применить функцию к каждому элементу
func Map[T, R any](slice []T, f func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = f(v)
    }
    return result
}

// Filter: оставить элементы, удовлетворяющие условию
func Filter[T any](slice []T, pred func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if pred(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce: свернуть в одно значение
func Reduce[T, R any](slice []T, init R, f func(R, T) R) R {
    acc := init
    for _, v := range slice {
        acc = f(acc, v)
    }
    return acc
}

// Использование:
nums := []int{1, 2, 3, 4, 5}

doubled := Map(nums, func(n int) int { return n * 2 })
// [2 4 6 8 10]

evens := Filter(nums, func(n int) bool { return n%2 == 0 })
// [2 4]

sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
// 15
```

---

## Когда использовать дженерики

**Используй дженерики, когда:**
- Логика одинакова для нескольких типов (контейнеры, утилиты для слайсов/мап)
- Хочешь типобезопасность без дублирования кода
- Пишешь библиотеку общего назначения

**Не используй дженерики, когда:**
- Достаточно интерфейса — он выражает поведение, дженерики выражают тип
- Метод нужен только для одного-двух конкретных типов
- Логика отличается для разных типов (это не параметризация, а полиморфизм)

```go
// НЕ нужны дженерики: io.Writer уже выражает нужную абстракцию
func WriteJSON(w io.Writer, v any) error {
    return json.NewEncoder(w).Encode(v)
}

// НУЖНЫ дженерики: одинаковая логика для любого упорядочиваемого типа
func Clamp[T cmp.Ordered](v, min, max T) T {
    if v < min { return min }
    if v > max { return max }
    return v
}
```

---

## Ограничения дженериков в Go

- Нет специализации: нельзя разное поведение для разных конкретных типов
- Нельзя использовать типовые параметры в методах (только в функциях и типах)
- Нельзя создать generic-метод на generic-типе с новым параметром

```go
type Container[T any] struct{ v T }

// НЕЛЬЗЯ: метод с новым типовым параметром
// func (c Container[T]) Convert[R any]() Container[R] { ... }  // ошибка компиляции

// МОЖНО: отдельная функция
func Convert[T, R any](c Container[T], f func(T) R) Container[R] {
    return Container[R]{v: f(c.v)}
}
```

---

## Итог

- `func F[T Constraint](v T)` — типовой параметр с ограничением
- `comparable` — поддерживает `==`; `cmp.Ordered` — поддерживает `<`, `>`
- `~T` в ограничении — включает типы с underlying type T
- Обобщённые типы: `type Stack[T any] struct{}`; методы используют `[T]`
- Дженерики для контейнеров и утилит; интерфейсы для поведения
