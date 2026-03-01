---
title: "Структуры"
description: "Пользовательские типы, вложенные структуры, теги, методы"
order: 5
---

# Структуры

Структура — составной тип, группирующий именованные поля разных типов.

## Объявление

```go
type User struct {
    Name  string
    Email string
    Age   int
}
```

## Создание экземпляра

```go
// Указываем все поля по порядку
u1 := User{"Алиса", "alice@example.com", 25}

// По именам полей (рекомендуется)
u2 := User{
    Name:  "Боб",
    Email: "bob@example.com",
    Age:   30,
}

// Нулевая структура
var u3 User // Name="", Email="", Age=0
```

## Доступ к полям

```go
fmt.Println(u2.Name)  // Боб
u2.Age = 31
```

## Указатели на структуры

```go
p := &User{Name: "Карл", Age: 28}
fmt.Println(p.Name) // Карл — автоматическое разыменование
```

Можно создать через `new`:

```go
p := new(User)  // *User, все поля нулевые
p.Name = "Дина"
```

## Анонимные структуры

Полезны для одноразовых структур:

```go
point := struct {
    X, Y int
}{10, 20}

fmt.Println(point.X) // 10
```

## Вложенные структуры

```go
type Address struct {
    City   string
    Street string
}

type Person struct {
    Name    string
    Address Address
}

p := Person{
    Name: "Иван",
    Address: Address{
        City:   "Москва",
        Street: "Тверская",
    },
}
fmt.Println(p.Address.City) // Москва
```

## Встраивание (embedding)

```go
type Person struct {
    Name string
    Address  // встраивание — без имени поля
}

p := Person{
    Name:    "Иван",
    Address: Address{City: "Москва", Street: "Тверская"},
}

// Доступ к полям напрямую:
fmt.Println(p.City) // Москва (вместо p.Address.City)
```

## Теги (tags)

Теги используются библиотеками для сериализации, валидации и ORM:

```go
type User struct {
    ID    int    `json:"id" db:"user_id"`
    Name  string `json:"name" validate:"required"`
    Email string `json:"email,omitempty"`
}
```

### Чтение тегов через reflect

```go
import "reflect"

t := reflect.TypeOf(User{})
field, _ := t.FieldByName("Email")
fmt.Println(field.Tag.Get("json")) // "email,omitempty"
```

## Сравнение структур

Структуры сравнимы, если все их поля сравнимы:

```go
a := User{Name: "Go", Age: 15}
b := User{Name: "Go", Age: 15}
fmt.Println(a == b) // true
```

> **Совет:** если структура содержит слайс или map, она не будет сравнимой через `==`. Используйте `reflect.DeepEqual`.

## Конструктор

Go не имеет конструкторов, но принято создавать функцию-фабрику:

```go
func NewUser(name, email string, age int) *User {
    return &User{
        Name:  name,
        Email: email,
        Age:   age,
    }
}
```
