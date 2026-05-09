---
title: "Структуры как коллекции"
description: "Slice of structs vs struct of slices, сортировка"
order: 5
---

# Структуры как коллекции

Слайсы структур — самый распространённый способ хранить набор объектов в Go. Разберём паттерны работы с такими коллекциями.

## Slice of Structs

Типичный паттерн:

```go
type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

users := []User{
    {ID: 1, Name: "Алиса", Email: "alice@example.com", Age: 30},
    {ID: 2, Name: "Боб", Email: "bob@example.com", Age: 25},
    {ID: 3, Name: "Карл", Email: "carl@example.com", Age: 35},
}
```

### Изменение элементов

Помни: элементы слайса возвращаются по значению в range. Для изменения — через индекс:

```go
// НЕВЕРНО: v — копия
for _, v := range users {
    v.Age++  // меняем копию
}
fmt.Println(users[0].Age)  // 30 — не изменился

// ВЕРНО: через индекс
for i := range users {
    users[i].Age++
}
fmt.Println(users[0].Age)  // 31

// АЛЬТЕРНАТИВА: slice of pointers
userPtrs := []*User{
    {ID: 1, Name: "Алиса", Age: 30},
}
for _, u := range userPtrs {
    u.Age++  // u — указатель, изменяем оригинал
}
fmt.Println(userPtrs[0].Age)  // 31
```

---

## Сортировка слайса структур

### sort.Slice

Самый гибкий способ:

```go
import "sort"

// Сортировка по возрасту (по возрастанию):
sort.Slice(users, func(i, j int) bool {
    return users[i].Age < users[j].Age
})

// По имени:
sort.Slice(users, func(i, j int) bool {
    return users[i].Name < users[j].Name
})

// По нескольким полям: по возрасту, затем по имени:
sort.Slice(users, func(i, j int) bool {
    if users[i].Age != users[j].Age {
        return users[i].Age < users[j].Age
    }
    return users[i].Name < users[j].Name
})
```

### sort.SliceStable — стабильная сортировка

Стабильная сортировка сохраняет исходный порядок элементов, которые считаются равными. Нужна когда данные уже частично упорядочены по другому полю и этот порядок важно не сломать.

Например, список пользователей отсортирован по имени. Нужно дополнительно отсортировать по возрасту, но среди одинакового возраста сохранить алфавитный порядок — `sort.Slice` этого не гарантирует, `sort.SliceStable` — гарантирует:

```go
// После этого: сначала по возрасту, среди равных — по имени (как было)
sort.SliceStable(users, func(i, j int) bool {
    return users[i].Age < users[j].Age
})
```

### Реализация sort.Interface

Для многократной сортировки одного типа — реализуй интерфейс `sort.Interface`:

```go
type ByAge []User

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

sort.Sort(ByAge(users))
```

> 💡 Интерфейс `sort.Interface` подробно разберём в главе 7.

---

> Сериализация структур в JSON и struct tags подробно разбираются в главе 6 «Указатели и структуры».

---

## Struct of Slices vs Slice of Structs

Два разных способа организации данных — выбор зависит от паттерна доступа.

### Slice of Structs (AoS — Array of Structs)

```go
type Point struct {
    X, Y, Z float64
    R, G, B uint8
}

points := []Point{...}  // каждый объект хранит все поля вместе
```

**Используй когда:**
- Работаешь с объектами как единым целым
- Часто передаёшь объект целиком
- Понятность кода важнее максимальной производительности

### Struct of Slices (SoA — Struct of Arrays)

```go
type PointCloud struct {
    X, Y, Z []float64  // отдельные слайсы для каждого поля
    R, G, B []uint8
}
```

**Используй когда:**
- Нужно обрабатывать только одно поле всех объектов (например, все X)
- Критична производительность (лучше cache locality для векторных операций)
- Работаешь с большими наборами данных в научных вычислениях

Для большинства бизнес-задач Slice of Structs понятнее и проще.

---

## Практический пример: фильтрация и сортировка

```go
package main

import (
    "fmt"
    "sort"
)

type Product struct {
    ID       int
    Name     string
    Price    float64
    Category string
}

func main() {
    products := []Product{
        {1, "Ноутбук", 79999, "electronics"},
        {2, "Книга", 590, "books"},
        {3, "Мышь", 1299, "electronics"},
        {4, "Роман", 399, "books"},
    }

    // Оставляем только электронику:
    var electronics []Product
    for _, p := range products {
        if p.Category == "electronics" {
            electronics = append(electronics, p)
        }
    }

    // Сортируем по цене:
    sort.Slice(electronics, func(i, j int) bool {
        return electronics[i].Price < electronics[j].Price
    })

    for _, p := range electronics {
        fmt.Printf("%s — %.0f₽\n", p.Name, p.Price)
    }
    // Мышь — 1299₽
    // Ноутбук — 79999₽
}
```

---

## Итог

- Слайс структур — основной паттерн для коллекций объектов
- Изменяй элементы через `s[i]`, а не через range-переменную (которая копия)
- `sort.Slice` — гибкая сортировка по любому критерию
- `sort.SliceStable` — когда важно не сломать существующий порядок среди равных элементов
- AoS (slice of structs) подходит для большинства задач, SoA (struct of slices) — для высоконагруженных векторных операций
