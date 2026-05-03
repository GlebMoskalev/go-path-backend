---
title: "Основы функций"
description: "Объявление, параметры, несколько возвращаемых значений и named return values"
order: 1
---

# Основы функций

Go позволяет возвращать несколько значений из функции — это меняет подход к обработке ошибок и делает код выразительнее. Разберём синтаксис от простого к сложному.

## Объявление функций

```go
func имяФункции(параметры) возвращаемыеТипы {
    // тело
}
```

```go
func greet(name string) string {
    return "Привет, " + name + "!"
}

func add(a, b int) int {
    return a + b
}

func noReturn() {
    fmt.Println("ничего не возвращаю")
}
```

### Сокращение одинаковых типов параметров

Если несколько параметров одного типа — тип можно указать только для последнего:

```go
// Полная форма:
func add(a int, b int, c int) int { return a + b + c }

// Сокращённая форма (идиоматично):
func add(a, b, c int) int { return a + b + c }

// Смешанные типы:
func format(prefix string, width, precision int) string {
    return fmt.Sprintf("%s: %*.*f", prefix, width, precision, 3.14)
}
```

---

## Несколько возвращаемых значений

Это одна из самых полезных возможностей Go. Функция может возвращать несколько значений:

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("деление на ноль")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 3)
    if err != nil {
        fmt.Println("Ошибка:", err)
        return
    }
    fmt.Printf("%.4f\n", result)  // 3.3333
}
```

Стандартный паттерн в Go: функция возвращает `(результат, error)`. Вызывающий код обязан проверить ошибку.

### Игнорирование значений через _

Если одно из возвращаемых значений не нужно:

```go
result, _ := divide(10, 3)  // игнорируем ошибку (не делай так в production!)

// Или первое значение:
_, err := os.Stat("file.txt")
if err != nil {
    fmt.Println("файл не существует")
}
```

### Функции как значения

```go
// Функция в переменной:
var fn func(int, int) int = add

// Вызов через переменную:
fmt.Println(fn(3, 4))  // 7

// Функция как параметр:
func apply(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, v := range nums {
        result[i] = f(v)
    }
    return result
}

doubled := apply([]int{1, 2, 3}, func(x int) int { return x * 2 })
fmt.Println(doubled)  // [2 4 6]
```

---

## Named return values — именованные возвращаемые значения

Возвращаемым значениям можно дать имена прямо в сигнатуре:

```go
func minmax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    for _, n := range nums[1:] {
        if n < min {
            min = n
        }
        if n > max {
            max = n
        }
    }
    return  // "голый" return — возвращает min и max
}

func main() {
    lo, hi := minmax([]int{3, 1, 4, 1, 5, 9, 2, 6})
    fmt.Println(lo, hi)  // 1 9
}
```

Именованные возвращаемые значения:
- Инициализируются нулевыми значениями при входе в функцию
- Видны как обычные локальные переменные в теле функции
- `return` без аргументов ("голый return") возвращает их текущие значения

### Когда использовать именованные возвращаемые значения

**Хорошо**: когда имена добавляют смысл к возвращаемым значениям:

```go
// Непонятно что такое (float64, float64):
func circleProps(radius float64) (float64, float64)

// Понятно:
func circleProps(radius float64) (area, circumference float64) {
    area = math.Pi * radius * radius
    circumference = 2 * math.Pi * radius
    return
}
```

**Хорошо**: в defer для изменения возвращаемых значений (паттерн с recover):

```go
func safeDiv(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("паника: %v", r)
        }
    }()
    result = a / b
    return
}
```

### Подводные камни именованных возвращаемых значений

**Ловушка 1**: `:=` внутри `if` создаёт новую локальную переменную вместо записи в именованное возвращаемое значение.

```go
// Плохо: user внутри if — это новая переменная, не именованный результат.
// После if именованный user по-прежнему nil, и голый return вернёт nil.
func getUser(id int) (user *User, err error) {
    if user, err := db.Find(id); err != nil {
        return nil, err
    }
    return // user == nil!
}

// Правильно: = вместо :=, присваиваем именованным переменным напрямую.
func getUser(id int) (user *User, err error) {
    user, err = db.Find(id)
    return
}
```

**Ловушка 2**: голый return в длинных функциях затрудняет чтение.

```go
// Плохо — непонятно что возвращается:
func complexCalc(x, y, z float64) (result float64, ok bool) {
    // ... 50 строк кода ...
    result = x * y * z
    ok = true
    return  // что именно возвращается?
}

// Лучше — явно:
func complexCalc(x, y, z float64) (float64, bool) {
    // ... 50 строк кода ...
    return x * y * z, true
}
```

Правило: **используй голые return только в коротких функциях**. В длинных — всегда явно указывай что возвращаешь.

---

## Практические примеры

### Функция разбора данных

```go
func parseCoord(s string) (lat, lon float64, err error) {
    parts := strings.Split(s, ",")
    if len(parts) != 2 {
        err = fmt.Errorf("неверный формат координат: %q", s)
        return
    }

    lat, err = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
    if err != nil {
        err = fmt.Errorf("неверная широта: %w", err)
        return
    }

    lon, err = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
    if err != nil {
        err = fmt.Errorf("неверная долгота: %w", err)
        return
    }

    return  // lat, lon, nil
}

func main() {
    lat, lon, err := parseCoord("55.7558, 37.6176")
    if err != nil {
        fmt.Println("Ошибка:", err)
        return
    }
    fmt.Printf("Широта: %.4f, Долгота: %.4f\n", lat, lon)
    // Широта: 55.7558, Долгота: 37.6176
}
```

---

## Типичные ошибки

**Ошибка 1**: Игнорировать ошибку из второго возвращаемого значения.

```go
result, _ := riskyOperation()  // тихая ошибка
result := riskyOperation()     // compile error если функция возвращает два значения
```

**Ошибка 2**: Изменить именованное возвращаемое через shadowing.

```go
func compute() (result int) {
    result = 10
    {
        result := 20  // новая переменная! именованный result не изменён
    }
    return  // вернёт 10, не 20
}
```

---

## Итог

- `func name(params) (returns)` — полная форма объявления
- Несколько типов одного вида: `func(a, b int)` вместо `func(a int, b int)`
- Несколько возвращаемых значений: идиоматичный паттерн `(result, error)`
- Именованные возвращаемые значения улучшают читаемость сигнатуры
- Голый `return` — только в коротких функциях, иначе путает
- Остерегайся shadowing именованных возвращаемых значений через `:=`
