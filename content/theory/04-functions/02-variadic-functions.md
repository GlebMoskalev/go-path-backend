---
title: "Вариативные функции"
description: "Функции с переменным числом аргументов: синтаксис ...T, передача слайса через ..., append как пример"
order: 2
---

# Вариативные функции

Вариативные функции принимают переменное количество аргументов одного типа. Самый знакомый пример — `fmt.Println`, которому можно передать сколько угодно значений.

## Синтаксис

```go
func имя(фиксированные параметры..., последний ...Тип) ВозвращаемыйТип {
    // последний параметр внутри функции — это []Тип
}
```

Три точки `...` перед типом последнего параметра:

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    fmt.Println(sum())            // 0
    fmt.Println(sum(1))           // 1
    fmt.Println(sum(1, 2, 3))     // 6
    fmt.Println(sum(1, 2, 3, 4, 5)) // 15
}
```

**Внутри функции** вариативний параметр `nums` — это обычный `[]int`. Все операции со слайсами применимы:

```go
func joinStrings(sep string, parts ...string) string {
    return strings.Join(parts, sep)
}

fmt.Println(joinStrings(", ", "яблоко", "банан", "вишня"))
// яблоко, банан, вишня
```

---

## Фиксированные параметры + variadic

Вариадический параметр всегда последний:

```go
func log(level string, args ...interface{}) {
    fmt.Printf("[%s] ", level)
    fmt.Println(args...)
}

log("INFO", "запуск сервера", "порт:", 8080)
// [INFO] запуск сервера порт: 8080
```

---

## Передача слайса в вариативную функцию

Если у тебя уже есть слайс и нужно передать его как variadic-аргументы — используй `...` при вызове:

```go
nums := []int{1, 2, 3, 4, 5}

// ОШИБКА: cannot use nums (type []int) as type int
fmt.Println(sum(nums))

// Правильно: распаковка слайса:
fmt.Println(sum(nums...))  // 15
```

Синтаксис `слайс...` говорит: «разбей слайс на отдельные аргументы».

**Важно**: при передаче слайса через `...` данные **не копируются**. Функция работает с тем же underlying array, поэтому изменения элементов видны снаружи:

```go
func setFirst(nums ...int) {
    if len(nums) > 0 {
        nums[0] = 999
    }
}

s := []int{1, 2, 3}
setFirst(s...)
fmt.Println(s) // [999 2 3] — оригинал изменился!
```

---

## append — классический пример

Встроенная функция `append` — самый используемый пример variadic в Go:

```go
func append(slice []Type, elems ...Type) []Type
```

```go
s := []int{1, 2, 3}

// Добавление одного элемента:
s = append(s, 4)

// Добавление нескольких:
s = append(s, 5, 6, 7)

// Добавление всех элементов другого слайса:
other := []int{8, 9, 10}
s = append(s, other...)  // распаковка через ...

fmt.Println(s)  // [1 2 3 4 5 6 7 8 9 10]
```

---

## Вариативные функции с интерфейсом any

`any` (псевдоним для `interface{}`) позволяет принимать аргументы любого типа — так устроены `fmt.Println`, `fmt.Printf` и другие:

```go
func printAll(sep string, values ...any) {
    parts := make([]string, 0, len(values))
    for _, v := range values {
        parts = append(parts, fmt.Sprintf("%v", v))
    }
    fmt.Println(strings.Join(parts, sep))
}

printAll(" | ", "Алиса", 30, true, 3.14)
// Алиса | 30 | true | 3.14
```

> 💡 Интерфейс `any` и работа с типами подробно разберём в главе 7 «Интерфейсы».

---

## Практический пример: функция логирования

```go
package main

import (
    "fmt"
    "time"
)

type LogLevel string

const (
    DEBUG LogLevel = "DEBUG"
    INFO  LogLevel = "INFO"
    ERROR LogLevel = "ERROR"
)

func logf(level LogLevel, format string, args ...any) {
    timestamp := time.Now().Format("15:04:05")
    message := fmt.Sprintf(format, args...)
    fmt.Printf("[%s] [%s] %s\n", timestamp, level, message)
}

func main() {
    logf(INFO, "Сервер запущен на порту %d", 8080)
    logf(DEBUG, "Подключение от %s", "192.168.1.1")
    logf(ERROR, "Не удалось открыть файл: %s", "config.yaml")
}
// [14:32:01] [INFO] Сервер запущен на порту 8080
// [14:32:01] [DEBUG] Подключение от 192.168.1.1
// [14:32:01] [ERROR] Не удалось открыть файл: config.yaml
```

---

## Типичные ошибки

**Ошибка 1**: Забыть `...` при передаче слайса.

```go
func maxOf(nums ...int) int { ... }

nums := []int{3, 1, 4, 1, 5}
maxOf(nums)    // ОШИБКА: cannot use nums ([]int) as int
maxOf(nums...) // OK
```

**Ошибка 2**: Ожидать что variadic параметр всегда не nil.

```go
func process(items ...string) {
    // items может быть nil если вызвали без аргументов!
    fmt.Println(len(items))  // 0, не паника — nil slice имеет длину 0
    for _, item := range items {  // безопасно
        fmt.Println(item)
    }
}

process()  // items == nil, но это нормально
```

**Ошибка 3**: Модификация variadic параметра — неожиданные побочные эффекты.

```go
func doubleSorted(nums ...int) []int {
    sort.Ints(nums)  // ВНИМАНИЕ: сортирует переданный слайс in-place!
    return nums
}

data := []int{3, 1, 2}
result := doubleSorted(data...)
fmt.Println(data)    // [1 2 3] — data изменился!
fmt.Println(result)  // [1 2 3]

// Безопасная версия — копируй:
func doubleSortedSafe(nums ...int) []int {
    sorted := make([]int, len(nums))
    copy(sorted, nums)
    sort.Ints(sorted)
    return sorted
}
```

---

## Итог

- `...T` — последний параметр принимает 0 или более значений типа T
- Внутри функции variadic параметр — обычный `[]T`
- `func(slice...)` — передать слайс как variadic-аргументы
- При передаче через `...` данные не копируются — будь осторожен с мутацией
- `append(s, other...)` — стандартный паттерн слияния слайсов
