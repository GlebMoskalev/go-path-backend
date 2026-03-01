---
title: "Функции"
description: "Объявление функций, множественные возвращаемые значения, вариативные параметры"
order: 1
---

# Функции

Функции — основной строительный блок программ на Go.

## Объявление

```go
func add(a int, b int) int {
    return a + b
}
```

Если параметры одного типа, тип можно указать один раз:

```go
func add(a, b int) int {
    return a + b
}
```

## Множественные возвращаемые значения

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("деление на ноль")
    }
    return a / b, nil
}

result, err := divide(10, 3)
```

## Именованные возвращаемые значения

```go
func swap(a, b int) (x, y int) {
    x = b
    y = a
    return // «голый» return возвращает x и y
}
```

> **Совет:** именованные возвращаемые значения полезны для документации, но `naked return` в длинных функциях снижает читаемость.

## Вариативные функции

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3)       // 6
sum(1, 2, 3, 4, 5) // 15
```

Передача слайса:

```go
nums := []int{1, 2, 3}
sum(nums...) // 6
```

## Функции как значения

```go
greet := func(name string) string {
    return "Привет, " + name + "!"
}

fmt.Println(greet("Гофер"))
```

## Функции как параметры

```go
func apply(nums []int, fn func(int) int) []int {
    result := make([]int, len(nums))
    for i, n := range nums {
        result[i] = fn(n)
    }
    return result
}

doubled := apply([]int{1, 2, 3}, func(n int) int {
    return n * 2
})
// [2, 4, 6]
```

## Функции как возвращаемое значение

```go
func multiplier(factor int) func(int) int {
    return func(n int) int {
        return n * factor
    }
}

double := multiplier(2)
triple := multiplier(3)

fmt.Println(double(5))  // 10
fmt.Println(triple(5))  // 15
```

## defer

`defer` откладывает вызов функции до выхода из текущей функции:

```go
func readFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close() // закроется при выходе из readFile

    // работа с файлом...
    return nil
}
```

Несколько `defer` выполняются в порядке LIFO (стек):

```go
defer fmt.Println("первый")
defer fmt.Println("второй")
defer fmt.Println("третий")
// Вывод: третий, второй, первый
```

## init()

Функция `init()` вызывается автоматически при загрузке пакета:

```go
func init() {
    // инициализация пакета
    fmt.Println("пакет загружен")
}
```

Можно иметь несколько `init()` в одном файле — они выполняются по порядку.
