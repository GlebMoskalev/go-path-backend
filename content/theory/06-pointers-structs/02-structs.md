---
title: "Структуры"
description: "Объявление, литералы, анонимные поля, сравнение структур, struct tags"
order: 2
---

# Структуры

Структура — основной способ создавать составные типы данных в Go. Нет классов, нет наследования — только структуры с методами.

## Объявление и инициализация

```go
type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}
```

### Способы создания

```go
// 1. Именованные поля (рекомендуется):
u1 := User{
    ID:    1,
    Name:  "Алиса",
    Email: "alice@example.com",
    Age:   30,
}

// 2. По порядку полей (избегай — хрупко при изменении структуры):
u2 := User{1, "Боб", "bob@example.com", 25}

// 3. Нулевое значение (все поля = нулевые значения):
var u3 User
fmt.Println(u3)  // {0  0}

// 4. Через new:
u4 := new(User)  // *User, все поля нулевые
u4.Name = "Карл"

// 5. Адрес литерала:
u5 := &User{Name: "Дима", Age: 28}
```

### Доступ к полям

```go
u := User{Name: "Алиса", Age: 30}

fmt.Println(u.Name)   // Алиса
u.Age = 31
fmt.Println(u.Age)    // 31

// Через указатель — Go автоматически разыменовывает:
p := &u
fmt.Println(p.Name)   // Алиса (p.Name == (*p).Name)
p.Age = 32            // изменяет оригинал
```

---

## Анонимные поля (встраивание)

Поле без имени — только тип. Его имя совпадает с именем типа:

```go
type Address struct {
    Street string
    City   string
}

type Person struct {
    Name    string
    Age     int
    Address  // анонимное поле — встраивание
}

p := Person{
    Name: "Алиса",
    Age:  30,
    Address: Address{
        Street: "Арбат, 1",
        City:   "Москва",
    },
}

// Доступ через имя встроенного типа:
fmt.Println(p.Address.City)   // Москва

// Или напрямую — "продвижение" полей:
fmt.Println(p.City)   // Москва (то же самое!)
fmt.Println(p.Street) // Арбат, 1
```

> 💡 Встраивание и продвижение методов подробно разберём в следующем уроке.

---

## Сравнение структур

Структуры сравниваемы через `==`, если все их поля сравниваемые:

```go
type Point struct{ X, Y int }

p1 := Point{1, 2}
p2 := Point{1, 2}
p3 := Point{1, 3}

fmt.Println(p1 == p2)  // true
fmt.Println(p1 == p3)  // false
```

Если структура содержит несравниваемые поля (слайс, map, функция) — `==` не работает:

```go
type BadStruct struct {
    Data []int  // слайс несравниваем
}

a := BadStruct{Data: []int{1, 2}}
b := BadStruct{Data: []int{1, 2}}
// a == b  // ОШИБКА компиляции: invalid operation: a == b (struct containing []int cannot be compared)

// Для сравнения слайсов — используй reflect.DeepEqual:
import "reflect"
fmt.Println(reflect.DeepEqual(a, b))  // true
```

---

## Struct tags

Теги — метаданные для полей структуры. Используются рефлексией и внешними пакетами:

```go
type Product struct {
    ID       int     `json:"id" db:"product_id"`
    Name     string  `json:"name" db:"name"`
    Price    float64 `json:"price" db:"price"`
    Internal string  `json:"-"`  // всегда пропускать
}
```

Формат тега: `` `ключ:"значение" ключ2:"значение2"` ``

Самые распространённые теги:

```go
type User struct {
    // json теги для encoding/json:
    ID    int    `json:"id"`
    Name  string `json:"name,omitempty"`  // пропустить если пустая строка
    Pass  string `json:"-"`               // никогда не сериализовывать

    // yaml теги для gopkg.in/yaml.v3:
    Config string `yaml:"config_path"`

    // db теги для sqlx/gorm:
    CreatedAt time.Time `db:"created_at" gorm:"column:created_at"`

    // validate теги для go-playground/validator:
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age"   validate:"min=0,max=150"`
}
```

### Чтение тегов через рефлексию

```go
import (
    "fmt"
    "reflect"
)

type Config struct {
    Host string `env:"HOST" default:"localhost"`
    Port int    `env:"PORT" default:"8080"`
}

func printTags(v interface{}) {
    t := reflect.TypeOf(v)
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fmt.Printf("Поле: %s, env: %s, default: %s\n",
            field.Name,
            field.Tag.Get("env"),
            field.Tag.Get("default"),
        )
    }
}

printTags(Config{})
// Поле: Host, env: HOST, default: localhost
// Поле: Port, env: PORT, default: 8080
```

> 💡 Рефлексию подробно разберём в главе 12 «Продвинутые темы».

---

## Структуры как значения

Как и массивы, структуры в Go — **значения**. Присваивание копирует всю структуру:

```go
type Point struct{ X, Y int }

p1 := Point{1, 2}
p2 := p1  // копия!

p2.X = 99
fmt.Println(p1)  // {1 2} — не изменился
fmt.Println(p2)  // {99 2}
```

Для разделения состояния — используй указатель:

```go
p3 := &p1  // указатель на p1
p3.X = 99
fmt.Println(p1)  // {99 2} — изменился!
```

---

## Анонимные структуры

Структуры без имени — полезны для одноразовых типов:

```go
// В коде:
point := struct{ X, Y int }{10, 20}
fmt.Println(point)  // {10 20}

// В тестах — часто используются для tabular tests:
tests := []struct {
    input    int
    expected int
}{
    {0, 0},
    {1, 1},
    {5, 120},
}

for _, tt := range tests {
    result := factorial(tt.input)
    if result != tt.expected {
        fmt.Printf("factorial(%d) = %d, want %d\n", tt.input, result, tt.expected)
    }
}

// Для JSON с динамической структурой:
response := struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}{200, "OK"}
```

---

## Итог

- Структуры — основа типов в Go; нет классов
- Всегда используй именованные поля при инициализации: `User{Name: "Alice"}`
- Структуры — значения: присваивание копирует
- Анонимные поля = встраивание; поля продвигаются
- Сравнение через `==` работает только если все поля сравниваемые
- Теги задают метаданные для полей; используются json, yaml, db и другими пакетами
