---
title: "Карты (map)"
description: "Хеш-таблицы в Go: создание, операции, перебор, вложенные карты"
order: 4
---

# Карты (map)

Map — ассоциативный массив (хеш-таблица), хранящий пары «ключ → значение».

## Создание

```go
// Литерал
m := map[string]int{
    "Go":     2009,
    "Rust":   2010,
    "Python": 1991,
}

// Через make
m := make(map[string]int)
```

## Операции

### Добавление / обновление

```go
m["Java"] = 1995
m["Go"] = 2012 // обновление
```

### Чтение

```go
year := m["Go"] // 2012
```

Если ключа нет, возвращается нулевое значение типа:

```go
fmt.Println(m["C++"]) // 0
```

### Проверка наличия ключа

```go
year, ok := m["Go"]
if ok {
    fmt.Println("Go:", year)
} else {
    fmt.Println("не найден")
}
```

Идиоматичный однострочник:

```go
if year, ok := m["Go"]; ok {
    fmt.Println(year)
}
```

### Удаление

```go
delete(m, "Python")
```

## Перебор

```go
for key, value := range m {
    fmt.Printf("%s → %d\n", key, value)
}
```

> **Важно:** порядок перебора map НЕ гарантирован. При каждом запуске он может быть разным.

## Длина

```go
fmt.Println(len(m)) // количество пар
```

## Nil-карта

```go
var m map[string]int // nil

// Чтение работает:
fmt.Println(m["key"]) // 0

// Запись вызывает panic:
m["key"] = 1 // panic: assignment to entry in nil map
```

Всегда инициализируйте карту перед записью!

## Map как множество (set)

```go
seen := make(map[string]bool)

words := []string{"go", "rust", "go", "python", "go"}
for _, w := range words {
    seen[w] = true
}

fmt.Println(seen)       // map[go:true python:true rust:true]
fmt.Println(seen["go"]) // true
```

Более эффективный вариант с `struct{}`:

```go
seen := make(map[string]struct{})

for _, w := range words {
    seen[w] = struct{}{}
}

if _, ok := seen["go"]; ok {
    fmt.Println("go уже есть")
}
```

## Вложенные карты

```go
graph := map[string]map[string]int{
    "A": {"B": 1, "C": 4},
    "B": {"C": 2},
}
```

> **Совет:** map в Go не потокобезопасен. Для конкурентного доступа используйте `sync.Map` или `sync.RWMutex`.
