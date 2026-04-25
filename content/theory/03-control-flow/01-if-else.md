---
title: "if и else"
description: "Инициализирующий оператор в if, вложенность, идиоматический Go с early return"
order: 1
---

# if и else

`if` в Go выглядит привычно, но имеет несколько отличий от Java/C/Python, которые делают код чище.

## Базовый синтаксис

```go
if условие {
    // ...
}

if условие {
    // ...
} else {
    // ...
}

if условие1 {
    // ...
} else if условие2 {
    // ...
} else {
    // ...
}
```

**Скобки вокруг условия — не нужны** (в отличие от C/Java, где они обязательны). Фигурные скобки — обязательны всегда, даже для одной строки:

```go
// ОШИБКА: нет фигурных скобок
if x > 0
    fmt.Println("positive")  // compile error

// Правильно:
if x > 0 {
    fmt.Println("positive")
}
```

---

## Инициализирующий оператор (init statement)

Это важная и уникальная особенность Go. В `if` можно включить короткое выражение, выполняемое перед проверкой условия:

```go
if инициализация; условие {
    // ...
}
```

Самый частый паттерн — проверка ошибки:

```go
// Без init statement:
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close()

// С init statement — компактнее:
if file, err := os.Open("data.txt"); err != nil {
    return err
} else {
    defer file.Close()
    // file доступен здесь...
}
// ...но не здесь! file виден только внутри if/else
```

**Ключевой момент**: переменные, объявленные в init statement, видны **только внутри блоков `if`/`else if`/`else`**, но не снаружи. Это уменьшает загрязнение scope.

Ещё примеры:

```go
// Чтение из map с проверкой наличия ключа:
if val, ok := userMap["alice"]; ok {
    fmt.Println("Нашли:", val)
} else {
    fmt.Println("Пользователь не найден")
}

// Type assertion с проверкой:
if str, ok := value.(string); ok {
    fmt.Println("Строка длиной", len(str))
}

// Выражение с результатом в условии:
if n := computeValue(); n > 100 {
    fmt.Println("Большое:", n)
} else if n > 0 {
    fmt.Println("Маленькое:", n)
} else {
    fmt.Println("Неположительное:", n)
}
```

---

## Early Return — идиоматический Go

В большинстве языков код часто выглядит как «пирамида смерти»: бесконечные вложения `if/else`. В Go принят противоположный подход — **early return**: проверяй ошибки и граничные случаи в начале функции, возвращай рано, оставляй основной путь без вложений.

**Плохо (pyramidal style):**

```go
func processUser(id int) error {
    user, err := getUser(id)
    if err == nil {
        if user.IsActive() {
            order, err := getOrder(user)
            if err == nil {
                if order.IsValid() {
                    return saveOrder(order)
                } else {
                    return errors.New("неверный заказ")
                }
            } else {
                return fmt.Errorf("заказ не найден: %w", err)
            }
        } else {
            return errors.New("пользователь неактивен")
        }
    } else {
        return fmt.Errorf("пользователь не найден: %w", err)
    }
}
```

**Хорошо (early return style):**

```go
func processUser(id int) error {
    user, err := getUser(id)
    if err != nil {
        return fmt.Errorf("пользователь не найден: %w", err)
    }

    if !user.IsActive() {
        return errors.New("пользователь неактивен")
    }

    order, err := getOrder(user)
    if err != nil {
        return fmt.Errorf("заказ не найден: %w", err)
    }

    if !order.IsValid() {
        return errors.New("неверный заказ")
    }

    return saveOrder(order)
}
```

Early return делает код «линейным» — читаешь сверху вниз, каждая строка одного уровня. Ошибки обрабатываются немедленно, а «счастливый путь» выделен отдельно.

---

## Практические паттерны

### Проверка входных данных

```go
func createUser(name, email string, age int) (*User, error) {
    if name == "" {
        return nil, errors.New("имя не может быть пустым")
    }
    if !isValidEmail(email) {
        return nil, fmt.Errorf("неверный формат email: %s", email)
    }
    if age < 0 || age > 150 {
        return nil, fmt.Errorf("недопустимый возраст: %d", age)
    }

    return &User{Name: name, Email: email, Age: age}, nil
}
```

### Условие с несколькими операндами

```go
if user != nil && user.IsActive() && !user.IsDeleted() {
    sendWelcomeEmail(user)
}
```

Короткое вычисление (`&&`) обеспечивает безопасность: если `user == nil`, остальные условия не проверяются.

---

## Типичные ошибки

**Ошибка 1**: Переменная из init statement не существует за пределами блока `if/else`.

```go
if err := doSomething(); err != nil {
    return err
}
// Здесь err недоступна!
fmt.Println(err)  // compile error: undefined: err
```

**Ошибка 2**: `else` после `return` — лишний.

```go
// Ненужный else после return:
if err != nil {
    return err
} else {  // этот else никогда не будет нужен
    process()
}

// Правильно (Go-style):
if err != nil {
    return err
}
process()
```

**Ошибка 3**: Забыть проверить второе возвращаемое значение функции.

```go
// Актуально для функций, возвращающих (значение, ошибка) или (значение, bool):
data, err := os.ReadFile("config.json")
// Если просто data := ... — ошибка игнорируется молча

if err != nil {
    return err
}
```

---

## Итог

- Скобки вокруг условия не нужны, фигурные скобки — обязательны
- Init statement в `if` уменьшает загрязнение scope: `if val, ok := m[k]; ok`
- Early return — идиоматический стиль: проверяй ошибки в начале, возвращай рано
- `else` после `return` — лишний, убирай
