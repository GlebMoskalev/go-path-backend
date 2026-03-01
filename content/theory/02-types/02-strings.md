---
title: "Строки"
description: "Строки в Go: UTF-8, руны, пакет strings, многострочные литералы"
order: 2
---

# Строки

Строки в Go — это неизменяемые последовательности байтов, закодированные в UTF-8.

## Объявление

```go
s1 := "Привет, мир!"
s2 := `Многострочный
литерал (raw string)`
```

Raw-строки (в обратных кавычках) не обрабатывают escape-последовательности:

```go
path := `C:\Users\gopher\docs`
```

## Длина

```go
s := "Hello"
fmt.Println(len(s))           // 5 — количество байтов

s2 := "Привет"
fmt.Println(len(s2))          // 12 — байтов (кириллица = 2 байта)
fmt.Println(utf8.RuneCountInString(s2)) // 6 — символов (рун)
```

## Руны

Руна (`rune`) — алиас для `int32`, представляет один символ Unicode.

```go
for i, r := range "Гофер" {
    fmt.Printf("байт %d: %c (U+%04X)\n", i, r, r)
}
```

## Конкатенация

```go
greeting := "Hello" + ", " + "World!"
```

Для множественной конкатенации используйте `strings.Builder`:

```go
var b strings.Builder
for i := 0; i < 100; i++ {
    b.WriteString("Go ")
}
result := b.String()
```

## Пакет strings

```go
import "strings"

s := "Hello, World!"

strings.Contains(s, "World")     // true
strings.HasPrefix(s, "Hello")    // true
strings.HasSuffix(s, "!")        // true
strings.ToUpper(s)               // "HELLO, WORLD!"
strings.ToLower(s)               // "hello, world!"
strings.TrimSpace("  Go  ")     // "Go"
strings.Replace(s, "World", "Go", 1) // "Hello, Go!"
strings.Split("a,b,c", ",")     // ["a", "b", "c"]
strings.Join([]string{"a","b"}, "-") // "a-b"
strings.Count(s, "l")           // 3
strings.Index(s, "World")       // 7
```

## Преобразование типов

```go
// строка ↔ []byte
b := []byte("hello")
s := string(b)

// строка ↔ []rune
r := []rune("Гофер")
s := string(r)

// число → строка
import "strconv"
s := strconv.Itoa(42)          // "42"
s := strconv.FormatFloat(3.14, 'f', 2, 64) // "3.14"

// строка → число
n, err := strconv.Atoi("42")
f, err := strconv.ParseFloat("3.14", 64)
```

## Сравнение строк

```go
fmt.Println("abc" == "abc")   // true
fmt.Println("abc" < "abd")    // true (лексикографически)
```

> **Совет:** строки в Go неизменяемы. Каждая операция создаёт новую строку. Для частых модификаций используйте `strings.Builder` или `[]byte`.
