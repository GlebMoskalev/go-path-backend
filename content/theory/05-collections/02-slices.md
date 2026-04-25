---
title: "Слайсы"
description: "Внутреннее устройство (ptr/len/cap), make, append с реаллокацией, copy, nil vs empty slice"
order: 2
---

# Слайсы

Слайс — самая важная структура данных в Go. Почти все коллекции в Go — это слайсы. Понимание их внутреннего устройства критично для написания эффективного кода.

## Внутреннее устройство

Слайс — это **дескриптор** из трёх полей, указывающий на участок массива:

```
Слайс:
┌─────────────┐
│   *array    │ → указатель на underlying array
├─────────────┤
│     len     │ — текущее количество элементов
├─────────────┤
│     cap     │ — максимум элементов без реаллокации
└─────────────┘
```

```go
a := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
s := a[2:5]  // слайс элементов 2, 3, 4

fmt.Println(s)       // [2 3 4]
fmt.Println(len(s))  // 3
fmt.Println(cap(s))  // 8 (от индекса 2 до конца массива)
```

```
Массив a:  [0][1][2][3][4][5][6][7][8][9]
                  ↑
Слайс s:         ptr   len=3   cap=8
                  └─────────────────────┘
```

**Слайс не хранит данные** — он только описывает диапазон существующего массива.

---

## Создание слайсов

### Литерал

```go
s := []int{1, 2, 3, 4, 5}
// автоматически создаётся underlying array и слайс поверх него
```

### make — создание с контролем len и cap

```go
s := make([]int, 5)       // len=5, cap=5, все элементы = 0
s := make([]int, 3, 10)   // len=3, cap=10
```

`make([]T, len, cap)` — стандартный способ создать слайс с предвыделенной памятью.

Когда знаешь примерный размер — указывай capacity:

```go
// Неэффективно: многократные реаллокации
var result []User
for _, id := range userIDs {
    user, _ := getUser(id)
    result = append(result, user)
}

// Эффективно: одна аллокация
result := make([]User, 0, len(userIDs))
for _, id := range userIDs {
    user, _ := getUser(id)
    result = append(result, user)
}
```

### Срез от массива или другого слайса

```go
a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

s1 := a[2:5]   // [2 3 4], len=3, cap=8
s2 := a[:3]    // [0 1 2], len=3, cap=10
s3 := a[7:]    // [7 8 9], len=3, cap=3
s4 := a[:]     // всё, len=10, cap=10

// Трёхсрезовая форма для явного управления cap:
s5 := a[2:5:6]  // [2 3 4], len=3, cap=4 (а не 8)
```

---

## append — добавление элементов

`append` добавляет элементы и возвращает новый слайс:

```go
s := []int{1, 2, 3}
s = append(s, 4)          // [1 2 3 4]
s = append(s, 5, 6, 7)    // [1 2 3 4 5 6 7]
```

**Всегда присваивай результат `append` обратно!**

### Реаллокация — когда cap исчерпан

Если `len == cap`, `append` создаёт новый underlying array (обычно в 2× больше), копирует данные и возвращает слайс с новым указателем:

```go
s := make([]int, 0, 3)
fmt.Printf("len=%d cap=%d %p\n", len(s), cap(s), s)

for i := 0; i < 6; i++ {
    s = append(s, i)
    fmt.Printf("len=%d cap=%d %p\n", len(s), cap(s), s)
}
// len=0 cap=3 0xc00001a060
// len=1 cap=3 0xc00001a060
// len=2 cap=3 0xc00001a060
// len=3 cap=3 0xc00001a060
// len=4 cap=6 0xc00003e060  ← новый адрес! реаллокация
// len=5 cap=6 0xc00003e060
// len=6 cap=6 0xc00003e060
```

**Последствия реаллокации**: слайсы, созданные до реаллокации, больше не разделяют underlying array с новым слайсом:

```go
original := make([]int, 3, 4)
copy(original, []int{1, 2, 3})

alias := original[:]  // тот же underlying array

original = append(original, 4)  // не реаллоцирует, cap=4 ещё есть
original[0] = 99
fmt.Println(alias[0])  // 99 — alias видит изменение!

original = append(original, 5)  // реаллоцирует! новый array
original[0] = 0
fmt.Println(alias[0])  // 99 — alias больше не связан с original
```

---

## copy — безопасное копирование

`copy(dst, src)` копирует элементы из src в dst, возвращает количество скопированных:

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, 3)

n := copy(dst, src)
fmt.Println(n)    // 3 — скопировал min(len(dst), len(src))
fmt.Println(dst)  // [1 2 3]

// Полная копия:
clone := make([]int, len(src))
copy(clone, src)
```

`copy` правильно обрабатывает **перекрывающиеся слайсы** (из одного underlying array):

```go
s := []int{1, 2, 3, 4, 5}
copy(s[1:], s[:])  // безопасный сдвиг вправо
fmt.Println(s)     // [1 1 2 3 4]
```

---

## nil slice vs empty slice

```go
var nilSlice []int      // nil slice: ptr=nil, len=0, cap=0
emptySlice := []int{}  // empty slice: ptr!=nil, len=0, cap=0
```

| Свойство | nil slice | empty slice |
|----------|-----------|-------------|
| `== nil` | `true` | `false` |
| `len()` | `0` | `0` |
| `cap()` | `0` | `0` |
| `append(s, ...)` | работает | работает |
| `for range` | безопасно | безопасно |
| JSON маршаллинг | `null` | `[]` |

```go
import "encoding/json"

var nilSlice []int
emptySlice := []int{}

n, _ := json.Marshal(nilSlice)
e, _ := json.Marshal(emptySlice)
fmt.Println(string(n))  // null
fmt.Println(string(e))  // []
```

**Когда важна разница**: при сериализации в JSON. Если API должен вернуть пустой массив, не nil — используй empty slice.

---

## Разделение underlying array — ловушка

Несколько слайсов могут разделять underlying array. Изменение через один влияет на другие:

```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:3]  // [2 3]

sub[0] = 99
fmt.Println(original)  // [1 99 3 4 5] — оригинал изменился!
fmt.Println(sub)       // [99 3]
```

Чтобы избежать неожиданных изменений — копируй:

```go
sub := make([]int, 2)
copy(sub, original[1:3])
sub[0] = 99
fmt.Println(original)  // [1 2 3 4 5] — не изменился
```

---

## Итог

- Слайс = ptr + len + cap; данные хранятся в underlying array
- `make([]T, len, cap)` — создание с предвыделенной памятью
- `append` может вызвать реаллокацию и разорвать связь с другими слайсами
- Всегда присваивай результат `append`: `s = append(s, ...)`
- `copy` — безопасное копирование, не разделяет память
- nil slice и empty slice ведут себя одинаково для append/range, но отличаются в JSON
- Срезы разделяют память — копируй, если нужна независимость
