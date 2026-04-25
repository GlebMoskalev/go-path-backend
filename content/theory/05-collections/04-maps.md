---
title: "Карты (Maps)"
description: "Hash map устройство, nil map ловушка, comma-ok, итерация, sync.Map"
order: 4
---

# Карты (Maps)

Map в Go — хэш-таблица с O(1) доступом по ключу. Это один из самых часто используемых типов данных.

## Объявление и создание

```go
// Объявление типа: map[KeyType]ValueType
var m map[string]int       // nil map — нельзя записывать!

m = make(map[string]int)   // инициализированная пустая map

// Литерал:
scores := map[string]int{
    "Алиса": 95,
    "Боб":   87,
    "Карл":  92,
}
```

---

## Nil map — классическая ловушка

**Чтение из nil map** — безопасно, возвращает нулевое значение:

```go
var m map[string]int
fmt.Println(m["alice"])  // 0 — нет паники
```

**Запись в nil map** — паника:

```go
var m map[string]int
m["alice"] = 95  // panic: assignment to entry in nil map!
```

Всегда инициализируй map перед записью:

```go
// Через make:
m := make(map[string]int)

// Через литерал (даже если пустой):
m := map[string]int{}
```

---

## Основные операции

```go
m := map[string]int{"a": 1, "b": 2}

// Чтение:
v := m["a"]      // 1
v = m["z"]       // 0 — ключ не существует, возвращает нулевое значение

// Запись:
m["c"] = 3

// Проверка существования — comma-ok:
val, ok := m["b"]
if ok {
    fmt.Println("найдено:", val)
} else {
    fmt.Println("ключ не существует")
}

// Удаление:
delete(m, "a")
fmt.Println(m)  // map[b:2 c:3]

// Размер:
fmt.Println(len(m))  // 2
```

### Comma-ok — правильная проверка ключа

Без comma-ok нельзя отличить "ключ с нулевым значением" от "ключ не существует":

```go
counters := map[string]int{
    "errors": 0,  // ноль ошибок
}

// ПЛОХО: невозможно различить отсутствие ключа и значение 0
v := counters["errors"]
if v == 0 {
    fmt.Println("нет ошибок или ключ не существует?")
}

// ХОРОШО: явная проверка наличия ключа
if v, ok := counters["errors"]; ok {
    fmt.Printf("ошибок: %d\n", v)  // ошибок: 0
} else {
    fmt.Println("счётчик ошибок не существует")
}
```

---

## Итерация по map

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

for k, v := range m {
    fmt.Printf("%s: %d\n", k, v)
}
// a: 1
// c: 3
// b: 2
// (порядок КАЖДЫЙ РАЗ разный!)
```

**Порядок итерации не определён** и намеренно рандомизируется. Никогда не полагайся на порядок.

Для детерминированного порядка — сортируй ключи:

```go
import "sort"

keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)

for _, k := range keys {
    fmt.Printf("%s: %d\n", k, m[k])
}
// a: 1
// b: 2
// c: 3
```

---

## Внутреннее устройство: почему порядок случаен

Map в Go — хэш-таблица с открытой адресацией. Ключи хэшируются и размещаются в бакетах. Начиная с Go 1, итерация начинается с **случайного бакета**, чтобы программисты не полагались на порядок — это намеренное дизайнерское решение.

Когда map превышает коэффициент загрузки (~6.5 элементов на бакет), происходит **рехэширование** — все элементы перераспределяются по новым бакетам. После этого порядок снова меняется.

---

## Вложенные map

```go
// Map of maps:
adjacency := map[string]map[string]int{
    "A": {"B": 1, "C": 4},
    "B": {"A": 1, "C": 2},
    "C": {"A": 4, "B": 2},
}

// Безопасный доступ к вложенной map:
if neighbors, ok := adjacency["A"]; ok {
    fmt.Println("соседи A:", neighbors)
}

// ЛОВУШКА: прямое добавление во вложенную nil map:
// adjacency["D"]["E"] = 5  // паника! adjacency["D"] == nil
adjacency["D"] = map[string]int{}  // сначала инициализируй
adjacency["D"]["E"] = 5
```

---

## Map как множество (set)

В Go нет встроенного типа Set. Эмулируется через `map[T]struct{}`:

```go
// struct{} — нулевой размер, не занимает память
seen := map[string]struct{}{}

words := []string{"apple", "banana", "apple", "cherry", "banana"}
for _, w := range words {
    seen[w] = struct{}{}
}

// Проверка вхождения:
if _, ok := seen["apple"]; ok {
    fmt.Println("apple есть в множестве")
}

// Все уникальные элементы:
fmt.Println(len(seen))  // 3
```

Альтернатива — `map[T]bool`, чуть менее идиоматично, но читается проще:

```go
seen := map[string]bool{}
seen["apple"] = true
if seen["apple"] {
    fmt.Println("apple есть")
}
```

---

## sync.Map для конкурентного доступа

Обычная `map` **не безопасна для конкурентного чтения/записи**:

```go
// ОПАСНО: race condition!
var m = map[string]int{}
go func() { m["key"] = 1 }()
go func() { fmt.Println(m["key"]) }()
```

Для конкурентного доступа есть два варианта:

### Вариант 1: обычная map + sync.RWMutex

```go
import "sync"

type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (sm *SafeMap) Set(key string, val int) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.m[key] = val
}

func (sm *SafeMap) Get(key string) (int, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    v, ok := sm.m[key]
    return v, ok
}
```

### Вариант 2: sync.Map

Оптимизирована для конкретных сценариев: много чтений и редкие записи, или разные ключи пишут разные горутины:

```go
import "sync"

var sm sync.Map

sm.Store("key", 42)

val, ok := sm.Load("key")
if ok {
    fmt.Println(val.(int))  // 42
}

// Атомарная операция "загрузи или сохрани":
actual, loaded := sm.LoadOrStore("key", 99)
fmt.Println(actual, loaded)  // 42 true — ключ уже был

sm.Delete("key")

// Итерация:
sm.Range(func(key, value any) bool {
    fmt.Println(key, value)
    return true  // вернуть false — остановить итерацию
})
```

> 💡 Конкурентность и sync-пакет подробно разберём в главе 10.

---

## Практический пример: подсчёт слов

```go
package main

import (
    "fmt"
    "sort"
    "strings"
)

func wordCount(text string) map[string]int {
    counts := make(map[string]int)
    for _, word := range strings.Fields(strings.ToLower(text)) {
        counts[word]++  // нулевое значение = 0, безопасно инкрементировать
    }
    return counts
}

func topN(counts map[string]int, n int) []string {
    type wordFreq struct {
        word  string
        count int
    }

    pairs := make([]wordFreq, 0, len(counts))
    for w, c := range counts {
        pairs = append(pairs, wordFreq{w, c})
    }
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].count > pairs[j].count
    })

    result := make([]string, 0, n)
    for i := 0; i < n && i < len(pairs); i++ {
        result = append(result, fmt.Sprintf("%s(%d)", pairs[i].word, pairs[i].count))
    }
    return result
}

func main() {
    text := "Go is expressive concise clean efficient Go is powerful"
    counts := wordCount(text)
    fmt.Println(topN(counts, 3))
    // [go(2) is(2) expressive(1)]
}
```

---

## Итог

- `map[K]V` — хэш-таблица, O(1) доступ
- Nil map: чтение безопасно (нулевое значение), запись — паника
- Всегда инициализируй: `make(map[K]V)` или литерал `map[K]V{}`
- Comma-ok для проверки: `val, ok := m[key]`
- Порядок итерации случаен — не полагайся на него
- `struct{}` как значение = множество (set) без расхода памяти
- Для конкурентного доступа: mutex или `sync.Map`
