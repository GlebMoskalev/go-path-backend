---
title: "Бенчмарки и примеры"
description: "testing.B, b.ResetTimer, b.ReportAllocs, Example-функции"
order: 3
---

# Бенчмарки и примеры

Go поддерживает два дополнительных вида тестов в том же `_test.go` файле: бенчмарки (измерение производительности) и примеры (исполняемая документация).

## Бенчмарки

Функция бенчмарка начинается с `Benchmark`, принимает `*testing.B`:

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

`b.N` — количество итераций. Инструмент тестирования автоматически подбирает `N` чтобы бенчмарк работал достаточно долго для получения точных результатов.

---

## Запуск бенчмарков

```bash
go test -bench=.               # все бенчмарки
go test -bench=BenchmarkAdd    # конкретный бенчмарк
go test -bench=. -benchtime=5s # запускать 5 секунд
go test -bench=. -count=3      # запустить 3 раза
go test -bench=. -benchmem     # показать аллокации памяти
```

Пример вывода:

```
BenchmarkAdd-8       1000000000    0.2843 ns/op
BenchmarkStringConcat-8   5000000    312 ns/op    48 B/op    3 allocs/op
```

- `-8` — количество ядер (GOMAXPROCS)
- `ns/op` — наносекунд на операцию
- `B/op` — байт аллоцировано на операцию
- `allocs/op` — количество аллокаций

---

## b.ResetTimer — исключить setup из измерений

```go
func BenchmarkProcessLargeSlice(b *testing.B) {
    // Дорогой setup: не должен входить в измерение
    data := make([]int, 1_000_000)
    for i := range data {
        data[i] = i
    }

    b.ResetTimer()  // начать отсчёт заново

    for i := 0; i < b.N; i++ {
        processSlice(data)
    }
}
```

### b.StopTimer / b.StartTimer

Для сложного setup внутри цикла:

```go
func BenchmarkSortRandom(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        data := generateRandomSlice(1000)  // setup не измеряем
        b.StartTimer()

        sort.Ints(data)  // только это измеряем
    }
}
```

---

## b.ReportAllocs — отчёт об аллокациях

`b.ReportAllocs()` эквивалентен флагу `-benchmem`, но на уровне конкретного бенчмарка:

```go
func BenchmarkStringConcat(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        s := ""
        for j := 0; j < 100; j++ {
            s += "x"  // много аллокаций
        }
        _ = s
    }
}
```

---

## Сравнение реализаций

Типичный паттерн: сравнить несколько подходов в одном файле:

```go
// strings_bench_test.go
package strings_test

import (
    "strings"
    "testing"
)

func BenchmarkConcatPlus(b *testing.B) {
    for i := 0; i < b.N; i++ {
        s := ""
        for j := 0; j < 100; j++ {
            s += "hello"
        }
        _ = s
    }
}

func BenchmarkConcatBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var sb strings.Builder
        for j := 0; j < 100; j++ {
            sb.WriteString("hello")
        }
        _ = sb.String()
    }
}

func BenchmarkConcatSlice(b *testing.B) {
    for i := 0; i < b.N; i++ {
        parts := make([]string, 100)
        for j := range parts {
            parts[j] = "hello"
        }
        _ = strings.Join(parts, "")
    }
}
```

```
BenchmarkConcatPlus-8        50000    24231 ns/op    25944 B/op    100 allocs/op
BenchmarkConcatBuilder-8   2000000      783 ns/op     4608 B/op      7 allocs/op
BenchmarkConcatSlice-8     1000000     1204 ns/op     5632 B/op      4 allocs/op
```

`strings.Builder` в 30 раз быстрее наивной конкатенации.

---

## Sub-бенчмарки

Аналог `t.Run` для бенчмарков:

```go
func BenchmarkProcess(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}

    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            data := make([]int, size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                process(data)
            }
        })
    }
}
```

Вывод:

```
BenchmarkProcess/size=10-8       5000000      312 ns/op
BenchmarkProcess/size=100-8       500000     3108 ns/op
BenchmarkProcess/size=1000-8       50000    31400 ns/op
BenchmarkProcess/size=10000-8       5000   312000 ns/op
```

Легко видно: сложность O(n).

---

## Example-функции

Функции `Example*` — это и тесты, и документация одновременно:

```go
func ExampleAdd() {
    fmt.Println(Add(2, 3))
    // Output: 5
}

func ExampleDivide() {
    result, err := Divide(10, 2)
    if err != nil {
        fmt.Println("ошибка:", err)
        return
    }
    fmt.Println(result)
    // Output: 5
}
```

Комментарий `// Output:` — ожидаемый вывод. Тест провалится если реальный вывод отличается.

```bash
go test -run "^Example" ./...
```

### Примеры в godoc

Example-функции отображаются в документации пакета. Именование:
- `ExampleAdd` — пример для функции `Add`
- `ExampleMyType_Method` — пример для метода `Method` типа `MyType`
- `Example` — пример для всего пакета
- `ExampleAdd_second` — второй пример для `Add` (суффикс через `_`)

```go
func ExampleAdd_negative() {
    fmt.Println(Add(-1, -2))
    // Output: -3
}
```

### Unordered output

Когда порядок вывода не детерминирован (например, итерация по map):

```go
func ExampleWordCount() {
    counts := WordCount("foo foo bar")
    for word, count := range counts {
        fmt.Printf("%s: %d\n", word, count)
    }
    // Unordered output:
    // foo: 2
    // bar: 1
}
```

---

## Итог

- `Benchmark*(*testing.B)` — функция бенчмарка; `b.N` подбирается автоматически
- `b.ResetTimer()` — исключить setup из измерений
- `b.ReportAllocs()` — показать аллокации памяти
- `-bench=. -benchmem` — запустить все бенчмарки с измерением памяти
- `Example*(...)` — пример-тест с `// Output:` для проверки и документации
- `// Unordered output:` — когда порядок вывода не гарантирован
