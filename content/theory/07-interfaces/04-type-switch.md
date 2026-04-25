---
title: "Type Switch"
description: "Type switch, паттерны использования, сравнение с type assertion"
order: 4
---

# Type Switch

Type switch — элегантный способ работать с несколькими возможными типами в одном блоке. Это специализированная форма switch для работы с интерфейсами.

## Синтаксис

```go
switch v := i.(type) {
case T1:
    // v имеет тип T1
case T2:
    // v имеет тип T2
case T3, T4:
    // v имеет тип interface{} (несколько типов в одном case)
default:
    // v имеет тот же тип что и i
}
```

Конструкция `i.(type)` — специальный синтаксис, работающий только внутри `switch`.

---

## Базовый пример

```go
func describe(i interface{}) string {
    switch v := i.(type) {
    case int:
        return fmt.Sprintf("целое: %d", v)
    case string:
        return fmt.Sprintf("строка %q длиной %d", v, len(v))
    case bool:
        if v {
            return "булево: true"
        }
        return "булево: false"
    case []int:
        return fmt.Sprintf("слайс int длиной %d", len(v))
    case nil:
        return "nil"
    default:
        return fmt.Sprintf("неизвестный тип %T", v)
    }
}

fmt.Println(describe(42))          // целое: 42
fmt.Println(describe("hello"))     // строка "hello" длиной 5
fmt.Println(describe(true))        // булево: true
fmt.Println(describe([]int{1,2}))  // слайс int длиной 2
fmt.Println(describe(nil))         // nil
fmt.Println(describe(3.14))        // неизвестный тип float64
```

---

## Несколько типов в одном case

Когда несколько типов обрабатываются одинаково:

```go
func isNumber(i interface{}) bool {
    switch i.(type) {
    case int, int8, int16, int32, int64,
         uint, uint8, uint16, uint32, uint64,
         float32, float64:
        return true
    }
    return false
}

fmt.Println(isNumber(42))    // true
fmt.Println(isNumber(3.14))  // true
fmt.Println(isNumber("hi"))  // false
```

**Внимание**: при нескольких типах в одном case переменная `v` имеет тип интерфейса (не конкретный тип):

```go
switch v := i.(type) {
case int, float64:
    // v тут — interface{}, не int и не float64
    // нужна дополнительная type assertion для использования
    fmt.Printf("%T: %v\n", v, v)
case string:
    // v тут — string, можно использовать напрямую
    fmt.Println(v + "!")
}
```

---

## Паттерны использования

### Обработка ошибок разных типов

```go
type NotFoundError struct{ Resource string }
func (e *NotFoundError) Error() string { return e.Resource + " not found" }

type ValidationError struct{ Field, Message string }
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func handleError(err error) {
    switch e := err.(type) {
    case *NotFoundError:
        fmt.Printf("Ресурс %s не найден, возвращаем 404\n", e.Resource)
    case *ValidationError:
        fmt.Printf("Ошибка валидации поля %s: %s\n", e.Field, e.Message)
    case nil:
        fmt.Println("Ошибок нет")
    default:
        fmt.Printf("Неизвестная ошибка: %v\n", e)
    }
}
```

### Десериализация JSON с динамической структурой

```go
func processJSONValue(v interface{}) string {
    switch val := v.(type) {
    case float64:
        return fmt.Sprintf("число: %g", val)
    case string:
        return fmt.Sprintf("строка: %q", val)
    case bool:
        return fmt.Sprintf("булево: %t", val)
    case []interface{}:
        return fmt.Sprintf("массив из %d элементов", len(val))
    case map[string]interface{}:
        return fmt.Sprintf("объект с %d ключами", len(val))
    case nil:
        return "null"
    default:
        return fmt.Sprintf("неизвестно: %T", val)
    }
}
```

### Visitor паттерн через type switch

```go
type Expr interface{ exprNode() }

type Number struct{ Val float64 }
type BinOp  struct{ Op string; Left, Right Expr }

func (n Number) exprNode() {}
func (b BinOp) exprNode()  {}

func evaluate(expr Expr) float64 {
    switch e := expr.(type) {
    case Number:
        return e.Val
    case BinOp:
        left := evaluate(e.Left)
        right := evaluate(e.Right)
        switch e.Op {
        case "+": return left + right
        case "-": return left - right
        case "*": return left * right
        case "/": return left / right
        }
    }
    panic("неизвестный тип выражения")
}

expr := BinOp{
    Op: "+",
    Left: Number{3},
    Right: BinOp{Op: "*", Left: Number{4}, Right: Number{5}},
}
fmt.Println(evaluate(expr))  // 23
```

---

## Type Switch vs Type Assertion

| | Type Switch | Type Assertion |
|---|---|---|
| Синтаксис | `switch v := i.(type)` | `v, ok := i.(T)` |
| Количество типов | Много | Один |
| Безопасность | Всегда безопасен | Нужна comma-ok форма |
| Когда использовать | 3+ возможных типа | 1-2 типа |

```go
// Для одного типа — type assertion проще:
if s, ok := i.(string); ok {
    fmt.Println(s)
}

// Для многих типов — type switch читаемее:
switch v := i.(type) {
case string:
    ...
case int:
    ...
case float64:
    ...
}
```

---

## Итог

- `switch v := i.(type)` — переключение по динамическому типу интерфейса
- В каждом case переменная `v` имеет конкретный тип (если один тип в case)
- При нескольких типах в case — `v` имеет тип интерфейса
- Используй для обработки ошибок, JSON-значений, AST-деревьев
- Для одного типа — предпочитай type assertion с comma-ok
