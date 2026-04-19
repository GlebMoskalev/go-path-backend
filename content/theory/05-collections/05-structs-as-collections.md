---
title: "Структуры как коллекции"
description: "Slice of structs vs struct of slices, сортировка, encoding/json"
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

```go
// Сохраняет относительный порядок равных элементов:
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

## encoding/json — сериализация

Пакет `encoding/json` знает как превратить структуры в JSON и обратно.

### Маршаллинг (struct → JSON)

```go
import "encoding/json"

type Product struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    InStock  bool    `json:"in_stock"`
}

p := Product{ID: 1, Name: "Ноутбук", Price: 79999.99, InStock: true}

data, err := json.Marshal(p)
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))
// {"id":1,"name":"Ноутбук","price":79999.99,"in_stock":true}
```

Теги `json:"..."` управляют именами ключей в JSON.

### Отступы для читаемости

```go
data, _ := json.MarshalIndent(p, "", "  ")
fmt.Println(string(data))
// {
//   "id": 1,
//   "name": "Ноутбук",
//   "price": 79999.99,
//   "in_stock": true
// }
```

### Анмаршаллинг (JSON → struct)

```go
jsonStr := `{"id":2,"name":"Мышь","price":1299.50,"in_stock":false}`

var product Product
err := json.Unmarshal([]byte(jsonStr), &product)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%+v\n", product)
// {ID:2 Name:Мышь Price:1299.5 InStock:false}
```

### Опции тегов

```go
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Password string `json:"-"`             // всегда пропускать в JSON
    Age      int    `json:"age,omitempty"` // пропустить если == 0
    Email    string `json:"email,omitempty"`
}

u := User{ID: 1, Name: "Алиса", Password: "secret", Age: 0}
data, _ := json.Marshal(u)
fmt.Println(string(data))
// {"id":1,"name":"Алиса"}  — Password и Age(omitempty) пропущены
```

### Маршаллинг слайса структур

```go
products := []Product{
    {ID: 1, Name: "Ноутбук", Price: 79999.99, InStock: true},
    {ID: 2, Name: "Мышь", Price: 1299.50, InStock: false},
}

data, _ := json.MarshalIndent(products, "", "  ")
fmt.Println(string(data))
// [
//   {
//     "id": 1,
//     "name": "Ноутбук",
//     ...
//   },
//   ...
// ]
```

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

## Практический пример: API-ответ

```go
package main

import (
    "encoding/json"
    "fmt"
    "sort"
)

type Product struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Category string  `json:"category"`
}

type APIResponse struct {
    Total    int       `json:"total"`
    Products []Product `json:"products"`
}

func main() {
    products := []Product{
        {1, "Ноутбук", 79999, "electronics"},
        {2, "Книга", 590, "books"},
        {3, "Мышь", 1299, "electronics"},
        {4, "Роман", 399, "books"},
    }

    // Сортируем по цене:
    sort.Slice(products, func(i, j int) bool {
        return products[i].Price < products[j].Price
    })

    response := APIResponse{
        Total:    len(products),
        Products: products,
    }

    data, _ := json.MarshalIndent(response, "", "  ")
    fmt.Println(string(data))
}
```

---

## Итог

- Слайс структур — основной паттерн для коллекций объектов
- Изменяй элементы через `s[i]`, а не через range-переменную (которая копия)
- `sort.Slice` — гибкая сортировка по любому критерию
- `json.Marshal` / `json.Unmarshal` — стандартная сериализация
- Теги `json:"name,omitempty"` и `json:"-"` управляют маршаллингом
- nil slice маршаллируется в `null`, empty slice — в `[]`
