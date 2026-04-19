---
title: "Рекурсия"
description: "Рекурсия, стек вызовов, хвостовая рекурсия (Go не оптимизирует!), мемоизация"
order: 5
---

# Рекурсия

Рекурсия — когда функция вызывает саму себя. Go полностью поддерживает рекурсию, но **не оптимизирует хвостовую рекурсию** — это важное отличие от некоторых других языков.

## Базовый пример

Классический пример — факториал:

```go
func factorial(n int) int {
    if n <= 1 {
        return 1  // базовый случай
    }
    return n * factorial(n-1)  // рекурсивный вызов
}

func main() {
    fmt.Println(factorial(5))   // 120
    fmt.Println(factorial(10))  // 3628800
}
```

Каждая рекурсивная функция имеет две части:
1. **Базовый случай** (base case) — условие выхода из рекурсии
2. **Рекурсивный шаг** — вызов себя с уменьшением задачи

Если забыть базовый случай — бесконечная рекурсия и stack overflow.

---

## Стек вызовов

Каждый вызов функции создаёт **stack frame** — блок памяти в стеке для локальных переменных и адреса возврата.

```
factorial(5)
  └── factorial(4)
        └── factorial(3)
              └── factorial(2)
                    └── factorial(1) → return 1
                  return 2 * 1 = 2
              return 3 * 2 = 6
        return 4 * 6 = 24
  return 5 * 24 = 120
```

Глубина рекурсии ограничена размером стека. В Go стек горутины начинается с 2-8 KB и **растёт динамически** до ~1 GB по умолчанию. Это лучше, чем в большинстве языков, но глубокая рекурсия всё равно может исчерпать стек:

```go
func inf(n int) int {
    return inf(n + 1)  // нет базового случая!
}
// runtime: goroutine stack exceeds 1000000000-byte limit
// fatal error: stack overflow
```

---

## Хвостовая рекурсия — Go НЕ оптимизирует

В функциональных языках (Haskell, Erlang) и некоторых других (Scala с TCO) компилятор оптимизирует **хвостовую рекурсию** — когда рекурсивный вызов последний в функции. Оптимизированная версия использует O(1) памяти стека.

**Go этого НЕ делает.** Каждый рекурсивный вызов создаёт новый stack frame.

```go
// Хвостовая рекурсия — последнее что делает функция = рекурсивный вызов:
func factorialTail(n, acc int) int {
    if n <= 1 {
        return acc
    }
    return factorialTail(n-1, n*acc)  // хвостовой вызов
}

// В Go это НЕ оптимизируется — стек растёт так же, как без TCO
// factorialTail(1000000, 1)  — stack overflow!
```

**Вывод**: для глубоких итераций в Go используй итеративный подход, а не рекурсию.

### Итеративный факториал

```go
func factorialIter(n int) int {
    result := 1
    for i := 2; i <= n; i++ {
        result *= i
    }
    return result
}
```

Итеративная версия использует O(1) памяти и работает для любых значений n.

---

## Практические применения рекурсии

Рекурсия уместна для задач, имеющих рекурсивную структуру по природе:

### Обход дерева

```go
type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

func sumTree(node *TreeNode) int {
    if node == nil {
        return 0
    }
    return node.Value + sumTree(node.Left) + sumTree(node.Right)
}

func printInOrder(node *TreeNode) {
    if node == nil {
        return
    }
    printInOrder(node.Left)
    fmt.Println(node.Value)
    printInOrder(node.Right)
}
```

### Обход файловой системы

```go
import (
    "fmt"
    "os"
    "path/filepath"
)

func walkDir(path string, depth int) error {
    entries, err := os.ReadDir(path)
    if err != nil {
        return err
    }

    indent := strings.Repeat("  ", depth)
    for _, entry := range entries {
        fmt.Printf("%s%s\n", indent, entry.Name())
        if entry.IsDir() {
            walkDir(filepath.Join(path, entry.Name()), depth+1)
        }
    }
    return nil
}
```

### Числа Фибоначчи (наивная рекурсия)

```go
func fib(n int) int {
    if n <= 1 {
        return n
    }
    return fib(n-1) + fib(n-2)
}
```

**Проблема**: экспоненциальная сложность O(2^n). `fib(40)` делает ~2 миллиарда вызовов.

---

## Мемоизация — кэшируй результаты

Мемоизация сохраняет уже вычисленные результаты, избегая повторных вычислений:

```go
package main

import "fmt"

func makeFib() func(int) int {
    cache := map[int]int{}

    var fib func(int) int
    fib = func(n int) int {
        if n <= 1 {
            return n
        }
        if cached, ok := cache[n]; ok {
            return cached  // возвращаем кэшированный результат
        }
        result := fib(n-1) + fib(n-2)
        cache[n] = result
        return result
    }

    return fib
}

func main() {
    fib := makeFib()
    fmt.Println(fib(10))  // 55
    fmt.Println(fib(50))  // 12586269025 — мгновенно
    fmt.Println(fib(90))  // 2880067194370816120 — мгновенно
}
```

С мемоизацией сложность падает с O(2^n) до O(n).

### Итеративное решение Фибоначчи — лучше для больших n

```go
func fibIter(n int) int {
    if n <= 1 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}

fmt.Println(fibIter(90))  // 2880067194370816120 — O(n) время, O(1) память
```

---

## Взаимная рекурсия

Две функции, вызывающие друг друга:

```go
func isEven(n int) bool {
    if n == 0 {
        return true
    }
    return isOdd(n - 1)
}

func isOdd(n int) bool {
    if n == 0 {
        return false
    }
    return isEven(n - 1)
}

// Для больших чисел — stack overflow!
// Лучше: return n%2 != 0
```

---

## Итог

- Рекурсия требует базового случая — без него stack overflow
- Go **не оптимизирует хвостовую рекурсию** — каждый вызов создаёт stack frame
- Для глубоких итераций используй итеративный подход
- Рекурсия уместна для задач с рекурсивной структурой: деревья, графы, файловая система
- Мемоизация превращает наивную экспоненциальную рекурсию в линейную
- Стек горутины в Go растёт динамически, но имеет предел (~1 GB)
