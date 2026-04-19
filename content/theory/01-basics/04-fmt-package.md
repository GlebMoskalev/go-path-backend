---
title: "Пакет fmt"
description: "Println, Printf, Sprintf, Fprintf, форматирование строк и чтение ввода"
order: 4
---

# Пакет fmt

Пакет `fmt` — это твой главный инструмент для вывода данных и чтения ввода. Умение правильно форматировать данные делает отладку быстрее, а код читаемее.

## Функции вывода

В пакете `fmt` три семейства функций. Разберём каждое.

### Print-семейство: куда выводить?

| Функция | Куда | Формат |
|---------|------|--------|
| `fmt.Print(...)` | stdout | без новой строки |
| `fmt.Println(...)` | stdout | пробелы между аргументами + `\n` |
| `fmt.Printf(format, ...)` | stdout | по шаблону |
| `fmt.Fprint(w, ...)` | `io.Writer` | без новой строки |
| `fmt.Fprintln(w, ...)` | `io.Writer` | пробелы + `\n` |
| `fmt.Fprintf(w, format, ...)` | `io.Writer` | по шаблону |
| `fmt.Sprint(...)` | строка | без новой строки |
| `fmt.Sprintln(...)` | строка | пробелы + `\n` |
| `fmt.Sprintf(format, ...)` | строка | по шаблону |

### Println — самый простой

```go
fmt.Println("Привет!")            // Привет!
fmt.Println("Имя:", "Алиса")     // Имя: Алиса
fmt.Println(1, 2, 3)             // 1 2 3
fmt.Println()                    // пустая строка
```

`Println` автоматически добавляет пробелы между аргументами разных типов и `\n` в конце.

### Printf — форматирование по шаблону

```go
name := "Алиса"
age := 30
score := 98.5

fmt.Printf("Имя: %s, возраст: %d\n", name, age)
// Имя: Алиса, возраст: 30

fmt.Printf("Результат: %.1f%%\n", score)
// Результат: 98.5%

fmt.Printf("%-10s | %5d\n", "Алиса", 30)  // выравнивание
//  Алиса      |    30
```

`Printf` **не добавляет** `\n` автоматически — не забывай писать явно.

### Sprintf — форматирование в строку

Наиболее полезен, когда нужно построить строку для дальнейшего использования:

```go
greeting := fmt.Sprintf("Добро пожаловать, %s! У тебя %d сообщений.", "Алиса", 5)
fmt.Println(greeting)
// Добро пожаловать, Алиса! У тебя 5 сообщений.
```

```go
// Форматирование ошибок
err := fmt.Errorf("пользователь %d не найден", 42)
fmt.Println(err)  // пользователь 42 не найден
```

### Fprintf — вывод в произвольный Writer

```go
import (
    "fmt"
    "os"
    "strings"
)

// Вывод в stderr (для ошибок):
fmt.Fprintf(os.Stderr, "ОШИБКА: %v\n", err)

// Вывод в строку через strings.Builder:
var sb strings.Builder
fmt.Fprintf(&sb, "Hello, %s!", "World")
fmt.Println(sb.String())  // Hello, World!
```

> 💡 `io.Writer` — интерфейс, который реализуют файлы, сетевые соединения, буферы и многое другое. Подробно в главе 12 «Продвинутые темы».

---

## Форматирующие глаголы (verbs)

### Универсальные

```go
type Point struct{ X, Y int }
p := Point{1, 2}

fmt.Printf("%v\n", p)   // {1 2}
fmt.Printf("%+v\n", p)  // {X:1 Y:2}  — с именами полей
fmt.Printf("%#v\n", p)  // main.Point{X:1, Y:2}  — Go-синтаксис
fmt.Printf("%T\n", p)   // main.Point
fmt.Printf("%%\n")       // %  — буквальный процент
```

### Числа

```go
n := 42

fmt.Printf("%d\n", n)   // 42       — десятичное
fmt.Printf("%b\n", n)   // 101010   — двоичное
fmt.Printf("%o\n", n)   // 52       — восьмеричное
fmt.Printf("%x\n", n)   // 2a       — шестнадцатеричное (строчные)
fmt.Printf("%X\n", n)   // 2A       — шестнадцатеричное (заглавные)
fmt.Printf("%c\n", 65)  // A        — символ Unicode

// Ширина и выравнивание:
fmt.Printf("%8d\n", 42)   //       42  — по правому краю, ширина 8
fmt.Printf("%-8d|\n", 42) // 42      |  — по левому краю
fmt.Printf("%08d\n", 42)  // 00000042  — с ведущими нулями
```

### Числа с плавающей точкой

```go
f := 3.14159265

fmt.Printf("%f\n", f)    // 3.141593   — стандартный (6 знаков после запятой)
fmt.Printf("%.2f\n", f)  // 3.14       — 2 знака после запятой
fmt.Printf("%e\n", f)    // 3.141593e+00  — научная нотация
fmt.Printf("%g\n", f)    // 3.14159265    — компактный формат
fmt.Printf("%9.3f\n", f) //     3.142  — ширина 9, 3 знака
```

### Строки

```go
s := "Привет"

fmt.Printf("%s\n", s)   // Привет
fmt.Printf("%q\n", s)   // "Привет"  — в кавычках, спецсимволы экранированы
fmt.Printf("%x\n", s)   // д0bfd180d0b8d0b2d0b5d18220  — байты в hex
fmt.Printf("%10s\n", s) //     Привет  — по правому краю
fmt.Printf("%-10s|\n", s) // Привет    |  — по левому краю
```

### Bool

```go
fmt.Printf("%t\n", true)   // true
fmt.Printf("%t\n", false)  // false
```

### Указатели

```go
x := 42
fmt.Printf("%p\n", &x)  // 0xc0000b4000  — адрес в памяти
```

---

## Практические примеры форматирования

### Таблица данных

```go
package main

import "fmt"

type Product struct {
    Name  string
    Price float64
    Count int
}

func main() {
    products := []Product{
        {"Ноутбук", 79999.99, 5},
        {"Мышь", 1299.50, 42},
        {"Клавиатура", 3499.00, 17},
    }

    fmt.Printf("%-15s %10s %6s\n", "Товар", "Цена", "Кол-во")
    fmt.Printf("%s\n", "─────────────────────────────────")
    for _, p := range products {
        fmt.Printf("%-15s %10.2f %6d\n", p.Name, p.Price, p.Count)
    }
}
// Товар            Цена Кол-во
// ─────────────────────────────────
// Ноутбук        79999.99      5
// Мышь            1299.50     42
// Клавиатура      3499.00     17
```

### Отладочный вывод

```go
func debugPrint(label string, v any) {
    fmt.Printf("[DEBUG] %s: %#v (type: %T)\n", label, v, v)
}

debugPrint("user", map[string]int{"alice": 30})
// [DEBUG] user: map[string]int{"alice":30} (type: map[string]int)
```

---

## Чтение ввода: Scan и Scanln

### Scan — читает разделённые пробелами значения

```go
package main

import "fmt"

func main() {
    var name string
    var age int

    fmt.Print("Введите имя и возраст: ")
    n, err := fmt.Scan(&name, &age)
    if err != nil {
        fmt.Println("Ошибка:", err)
        return
    }
    fmt.Printf("Считано %d значений: %s, %d лет\n", n, name, age)
}
```

`Scan` читает разделённые пробелами значения и разносит по переменным. Обязательно передавай **указатели** (`&name`, `&age`).

### Scanln — чтение до конца строки

```go
var first, last string
fmt.Print("Имя Фамилия: ")
fmt.Scanln(&first, &last)
fmt.Printf("Привет, %s %s!\n", first, last)
```

`Scanln` останавливается на `\n`, а не на EOF.

### Scanf — ввод по шаблону

```go
var day, month, year int
fmt.Print("Введите дату (дд.мм.гггг): ")
fmt.Scanf("%d.%d.%d", &day, &month, &year)
fmt.Printf("День: %d, Месяц: %d, Год: %d\n", day, month, year)
```

**Практическое замечание**: для серьёзных программ предпочитай `bufio.Scanner` для чтения строк — он надёжнее обрабатывает Unicode и длинные строки:

```go
import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("Введите строку: ")
    if scanner.Scan() {
        line := scanner.Text()
        fmt.Println("Вы ввели:", line)
    }
}
```

---

## Типичные ошибки

**Ошибка 1**: Неверное количество аргументов в Printf.

```go
fmt.Printf("%s %d", "Алиса")
// Алиса %!d(MISSING)  — Go не падает, но предупреждает
```

**Ошибка 2**: Забыть `\n` в Printf.

```go
fmt.Printf("Результат: %d", 42)  // нет переноса строки
fmt.Printf("Следующая строка")    // на той же строке!
// Результат: 42Следующая строка
```

**Ошибка 3**: Передать значение вместо указателя в Scan.

```go
var x int
fmt.Scan(x)   // ОШИБКА: передано значение, не указатель
fmt.Scan(&x)  // OK
```

---

## Итог

- `Println` — для быстрого вывода нескольких значений
- `Printf` — для форматированного вывода по шаблону, не добавляет `\n`
- `Sprintf` — для построения форматированной строки (не вывода)
- `Fprintf` — для записи в произвольный `io.Writer` (файл, stderr, буфер)
- `%v` — универсальный глагол, `%T` — тип, `%#v` — Go-синтаксис
- В `Scan`-функциях всегда передавай указатели
