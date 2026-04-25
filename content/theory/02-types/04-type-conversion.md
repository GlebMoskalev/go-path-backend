---
title: "Преобразование типов"
description: "Явное приведение типов, пакет strconv, основные функции парсинга и форматирования"
order: 4
---

# Преобразование типов

В Go нет неявного преобразования типов. Никогда. Это одно из ключевых дизайнерских решений языка, которое поначалу кажется многословным, но на практике предотвращает целые классы ошибок.

## Явное преобразование

Синтаксис: `T(выражение)` где `T` — целевой тип.

```go
var i int = 42
var f float64 = float64(i)  // int → float64
var u uint = uint(f)        // float64 → uint

fmt.Println(i, f, u)  // 42 42 42
```

Примеры типичных преобразований:

```go
// Числовые типы
var x int32 = 100
var y int64 = int64(x)
var z float64 = float64(x)

// int → float для правильного деления
a, b := 7, 2
ratio := float64(a) / float64(b)
fmt.Println(ratio)  // 3.5 (а не 3!)

// float → int (дробная часть ОТРЕЗАЕТСЯ, не округляется)
pi := 3.99
n := int(pi)
fmt.Println(n)  // 3, не 4!

// byte ↔ rune ↔ int ↔ string
var r rune = 'А'           // Unicode code point для 'А'
var b byte = 65             // ASCII код для 'A'
fmt.Println(string(r))     // А
fmt.Println(string(b))     // A
fmt.Println(string(72))    // H (код буквы H в Unicode)

// []byte ↔ string (копирует данные)
s := "Hello"
bytes := []byte(s)
bytes[0] = 'h'
fmt.Println(string(bytes))  // hello (s не изменилась!)
```

### Когда преобразование возможно

Числовые типы преобразуются свободно (с возможной потерей данных). Нельзя преобразовывать несовместимые типы:

```go
type Celsius float64
type Fahrenheit float64

c := Celsius(100)
// f := Fahrenheit(c)  // ОШИБКА: нельзя преобразовать Celsius в Fahrenheit напрямую
f := Fahrenheit(float64(c)*9/5 + 32)  // OK: через базовый тип
fmt.Printf("%.1f°C = %.1f°F\n", c, f)  // 100.0°C = 212.0°F
```

---

## Потеря данных при преобразовании

Преобразование может привести к потере точности или переполнению:

```go
// Потеря дробной части
f := 3.99
i := int(f)
fmt.Println(i)  // 3 — НЕ 4

// Переполнение при сужении типа
big := int64(1000)
small := int8(big)  // 1000 не влезает в int8 (макс 127)
fmt.Println(small)  // -24 — тихое переполнение!

// Потеря точности float64 → float32
precise := float64(3.141592653589793)
approx := float32(precise)
fmt.Printf("%.15f\n", precise)           // 3.141592653589793
fmt.Printf("%.15f\n", float64(approx))  // 3.141592741012573 — потеря точности
```

---

## Пакет strconv — конвертация строк

Для преобразований между строками и числами используй `strconv`. `fmt.Sprintf` тоже работает, но `strconv` быстрее и правильнее для этих задач.

### Строка → число

```go
import "strconv"

// Atoi — string в int (самая частая операция)
n, err := strconv.Atoi("42")
if err != nil {
    fmt.Println("Ошибка:", err)
    return
}
fmt.Println(n + 1)  // 43

// Ошибка при неверном вводе:
_, err = strconv.Atoi("abc")
fmt.Println(err)  // strconv.Atoi: parsing "abc": invalid syntax

// ParseInt — гибкая версия с базой и битовым размером
i64, err := strconv.ParseInt("-128", 10, 8)  // base 10, int8
fmt.Println(i64, err)  // -128 <nil>

i64, err = strconv.ParseInt("FF", 16, 16)   // шестнадцатеричное число
fmt.Println(i64, err)  // 255 <nil>

// ParseFloat
f, err := strconv.ParseFloat("3.14159", 64)  // 64-битная точность
fmt.Println(f, err)  // 3.14159 <nil>

// ParseBool
b, err := strconv.ParseBool("true")
fmt.Println(b, err)  // true <nil>

b, _ = strconv.ParseBool("1")
fmt.Println(b)  // true

b, _ = strconv.ParseBool("T")
fmt.Println(b)  // true

b, _ = strconv.ParseBool("0")
fmt.Println(b)  // false
```

### Число → строка

```go
// Itoa — int в string
s := strconv.Itoa(42)
fmt.Printf("%q\n", s)  // "42"

// FormatInt с указанием основания
hex := strconv.FormatInt(255, 16)
fmt.Println(hex)  // ff

bin := strconv.FormatInt(42, 2)
fmt.Println(bin)  // 101010

// FormatFloat
pi := strconv.FormatFloat(3.14159265, 'f', 2, 64)
fmt.Println(pi)  // 3.14

sci := strconv.FormatFloat(0.0001234, 'e', 3, 64)
fmt.Println(sci)  // 1.234e-04

// 'g' — автоматически выбирает f или e:
compact := strconv.FormatFloat(3.14159265, 'g', -1, 64)
fmt.Println(compact)  // 3.14159265

// FormatBool
fmt.Println(strconv.FormatBool(true))   // true
fmt.Println(strconv.FormatBool(false))  // false
```

### Практический пример: парсинг конфига

```go
package main

import (
    "fmt"
    "strconv"
)

func parseConfig(params map[string]string) error {
    portStr := params["port"]
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return fmt.Errorf("неверный порт %q: %w", portStr, err)
    }
    if port < 1 || port > 65535 {
        return fmt.Errorf("порт %d вне допустимого диапазона", port)
    }

    debugStr := params["debug"]
    debug, err := strconv.ParseBool(debugStr)
    if err != nil {
        return fmt.Errorf("неверное значение debug %q: %w", debugStr, err)
    }

    timeoutStr := params["timeout_seconds"]
    timeout, err := strconv.ParseFloat(timeoutStr, 64)
    if err != nil {
        return fmt.Errorf("неверный timeout %q: %w", timeoutStr, err)
    }

    fmt.Printf("Порт: %d, Отладка: %v, Таймаут: %.1f с\n", port, debug, timeout)
    return nil
}

func main() {
    config := map[string]string{
        "port":             "8080",
        "debug":            "true",
        "timeout_seconds":  "30.5",
    }

    if err := parseConfig(config); err != nil {
        fmt.Println("Ошибка:", err)
    }
    // Порт: 8080, Отладка: true, Таймаут: 30.5 с
}
```

---

## fmt.Sprintf vs strconv

Оба способа работают, но для чистого преобразования числа в строку `strconv` быстрее:

```go
n := 42

// fmt.Sprintf — удобно, но медленнее (аллоцирует строку + форматирует)
s1 := fmt.Sprintf("%d", n)

// strconv.Itoa — быстрее, специализировано для этой задачи
s2 := strconv.Itoa(n)

fmt.Println(s1 == s2)  // true — результат одинаков
```

Используй `fmt.Sprintf` для сложного форматирования ("пользователь %s (id: %d)"), `strconv` — для чистой конвертации числа в строку.

---

## Упоминание unsafe

Пакет `unsafe` позволяет обходить систему типов Go, включая конвертации без копирования:

```go
import "unsafe"

s := "hello"
// Конвертация string в []byte без копирования (ОПАСНО):
b := unsafe.Slice(unsafe.StringData(s), len(s))
```

**Не используй `unsafe` без острой необходимости.** Это инструмент для системного программирования, нарушает гарантии безопасности Go, и код с `unsafe` может сломаться при обновлении версии Go.

---

## Итог

- Go требует **явных** преобразований типов — никакой магии
- `T(x)` — синтаксис преобразования; дробная часть при float→int отрезается
- Сужение типа может привести к тихому переполнению
- `strconv` — стандартный способ парсинга строк в числа и обратно
- `strconv.Atoi` / `strconv.Itoa` — для int; `ParseFloat` / `FormatFloat` — для float
- Всегда проверяй ошибку при парсинге пользовательского ввода
