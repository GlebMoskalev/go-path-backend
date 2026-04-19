---
title: "Трюки со слайсами"
description: "Удаление элемента, вставка, reverse, фильтрация без лишних аллокаций"
order: 3
---

# Трюки со слайсами

Несколько стандартных операций со слайсами, которые часто нужны на практике, но не встроены в язык.

## Удаление элемента

### Удаление с сохранением порядка

```go
func remove(s []int, i int) []int {
    return append(s[:i], s[i+1:]...)
}

s := []int{1, 2, 3, 4, 5}
s = remove(s, 2)   // удаляем элемент с индексом 2 (значение 3)
fmt.Println(s)     // [1 2 4 5]
```

**Как работает**: `append(s[:i], s[i+1:]...)` — объединяем левую и правую части, `s[i+1:]` разворачивается как variadic.

**Предупреждение**: `s[:i]` и `s[i+1:]` разделяют underlying array! После операции `s[i+1:]` перезаписывает элемент в позиции `i`. Исходный слайс **изменяется in-place**.

```go
s := []int{1, 2, 3, 4, 5}
ref := s[2]  // 3
_ = ref
s = remove(s, 1)  // удаляем 2
fmt.Println(s)    // [1 3 4 5]
fmt.Println(s[:5]) // panic! len уменьшился, но элемент 5 всё ещё в памяти
```

### Удаление без сохранения порядка (быстрее)

Если порядок не важен — заменяем удаляемый элемент последним и уменьшаем длину:

```go
func removeUnordered(s []int, i int) []int {
    s[i] = s[len(s)-1]  // заменяем удаляемый последним
    return s[:len(s)-1]  // уменьшаем длину
}

s := []int{1, 2, 3, 4, 5}
s = removeUnordered(s, 1)  // удаляем элемент с индексом 1
fmt.Println(s)  // [1 5 3 4] — порядок нарушен, но O(1)
```

Это O(1) против O(n) для сохраняющего порядок варианта.

---

## Вставка элемента

```go
func insert(s []int, i int, val int) []int {
    s = append(s, 0)       // увеличиваем длину
    copy(s[i+1:], s[i:])  // сдвигаем элементы вправо
    s[i] = val
    return s
}

s := []int{1, 2, 4, 5}
s = insert(s, 2, 3)  // вставляем 3 перед индексом 2
fmt.Println(s)        // [1 2 3 4 5]
```

Или через append:

```go
func insertAppend(s []int, i int, val int) []int {
    return append(s[:i], append([]int{val}, s[i:]...)...)
}
```

Вторая версия создаёт временный слайс — менее эффективна. Первая работает in-place.

---

## Reverse — разворот слайса

```go
func reverse(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]  // swap
    }
}

s := []int{1, 2, 3, 4, 5}
reverse(s)
fmt.Println(s)  // [5 4 3 2 1]
```

Разворот in-place, O(n) время, O(1) память.

Обобщённая версия через дженерики (Go 1.18+):

```go
func reverseGeneric[T any](s []T) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}
```

---

## Фильтрация без выделения памяти

Стандартный подход — создать новый слайс. Но можно фильтровать in-place, переиспользуя underlying array:

```go
// Фильтрация с аллокацией нового слайса:
func filterNew(s []int, keep func(int) bool) []int {
    result := make([]int, 0, len(s))
    for _, v := range s {
        if keep(v) {
            result = append(result, v)
        }
    }
    return result
}

// Фильтрация in-place (без аллокации):
func filterInPlace(s []int, keep func(int) bool) []int {
    n := 0
    for _, v := range s {
        if keep(v) {
            s[n] = v
            n++
        }
    }
    return s[:n]
}

nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
evens := filterInPlace(nums, func(x int) bool { return x%2 == 0 })
fmt.Println(evens)  // [2 4 6 8 10]
```

**Предупреждение**: `filterInPlace` модифицирует оригинальный слайс. Если нужно сохранить оригинал — используй `filterNew`.

---

## Дедупликация (удаление дубликатов)

```go
// Для отсортированного слайса — O(n):
func deduplicateSorted(s []int) []int {
    if len(s) == 0 {
        return s
    }
    n := 1
    for i := 1; i < len(s); i++ {
        if s[i] != s[i-1] {
            s[n] = s[i]
            n++
        }
    }
    return s[:n]
}

// Для неотсортированного — O(n) с map:
func deduplicate(s []int) []int {
    seen := make(map[int]struct{}, len(s))
    result := make([]int, 0, len(s))
    for _, v := range s {
        if _, ok := seen[v]; !ok {
            seen[v] = struct{}{}
            result = append(result, v)
        }
    }
    return result
}
```

---

## Поиск элемента

```go
// Линейный поиск:
func contains(s []int, target int) bool {
    for _, v := range s {
        if v == target {
            return true
        }
    }
    return false
}

// Бинарный поиск (для отсортированного слайса):
import "sort"

s := []int{1, 3, 5, 7, 9, 11}
i := sort.SearchInts(s, 7)
if i < len(s) && s[i] == 7 {
    fmt.Println("найден на позиции", i)  // найден на позиции 3
}
```

---

## Разбивка на чанки

```go
func chunks(s []int, size int) [][]int {
    var result [][]int
    for len(s) > 0 {
        if len(s) < size {
            size = len(s)
        }
        result = append(result, s[:size])
        s = s[size:]
    }
    return result
}

nums := []int{1, 2, 3, 4, 5, 6, 7}
for _, chunk := range chunks(nums, 3) {
    fmt.Println(chunk)
}
// [1 2 3]
// [4 5 6]
// [7]
```

---

## Сортировка

Пакет `sort` содержит функции для сортировки слайсов:

```go
import "sort"

nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
sort.Ints(nums)
fmt.Println(nums)  // [1 1 2 3 4 5 6 9]

words := []string{"банан", "яблоко", "вишня"}
sort.Strings(words)
fmt.Println(words)  // [банан вишня яблоко]

// Кастомная сортировка через sort.Slice:
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Боб", 30},
    {"Алиса", 25},
    {"Карл", 35},
}

sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
fmt.Println(people)
// [{Алиса 25} {Боб 30} {Карл 35}]
```

---

## Итог

| Операция | Метод | Сложность |
|----------|-------|-----------|
| Удаление (с порядком) | `append(s[:i], s[i+1:]...)` | O(n) |
| Удаление (без порядка) | замена последним + сужение | O(1) |
| Вставка | сдвиг + присваивание | O(n) |
| Разворот | swap с двух концов | O(n) |
| Фильтрация in-place | перезапись + сужение | O(n) без аллокации |
| Сортировка | `sort.Slice` | O(n log n) |

Большинство трюков работают in-place — переиспользуют underlying array, избегая лишних аллокаций.
