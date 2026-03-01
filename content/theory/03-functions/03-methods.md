---
title: "Методы"
description: "Методы на типах, ресиверы по значению и по указателю"
order: 3
---

# Методы

Метод — это функция, привязанная к типу через ресивер (receiver).

## Объявление метода

```go
type Rect struct {
    Width, Height float64
}

func (r Rect) Area() float64 {
    return r.Width * r.Height
}

func (r Rect) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}
```

## Вызов

```go
rect := Rect{Width: 10, Height: 5}
fmt.Println(rect.Area())      // 50
fmt.Println(rect.Perimeter()) // 30
```

## Ресивер по значению vs по указателю

### По значению

Метод получает **копию** структуры. Изменения не влияют на оригинал:

```go
func (r Rect) Scale(factor float64) Rect {
    return Rect{
        Width:  r.Width * factor,
        Height: r.Height * factor,
    }
}
```

### По указателю

Метод получает **ссылку**. Может изменять оригинал:

```go
func (r *Rect) ScaleInPlace(factor float64) {
    r.Width *= factor
    r.Height *= factor
}
```

```go
rect := Rect{Width: 10, Height: 5}
rect.ScaleInPlace(2)
fmt.Println(rect) // {20 10}
```

## Когда использовать указатель

Используйте ресивер-указатель, если:

1. Метод **изменяет** ресивер
2. Структура **большая** (избегаем копирования)
3. Для **единообразия** — если хоть один метод использует указатель, все методы должны

## Методы на любых типах

Можно определять методы на любых типах, объявленных в том же пакете:

```go
type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
    return float64(c)*9/5 + 32
}

temp := Celsius(100)
fmt.Println(temp.ToFahrenheit()) // 212
```

```go
type StringSlice []string

func (ss StringSlice) Contains(target string) bool {
    for _, s := range ss {
        if s == target {
            return true
        }
    }
    return false
}
```

## Метод как значение

```go
rect := Rect{Width: 10, Height: 5}
areaFn := rect.Area  // привязан к конкретному rect
fmt.Println(areaFn()) // 50
```

## Метод-выражение

```go
areaFn := Rect.Area
fmt.Println(areaFn(rect)) // 50 — rect передаётся первым аргументом
```

## Встраивание и продвижение методов

```go
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return a.Name + " говорит!"
}

type Dog struct {
    Animal  // встраивание
    Breed string
}

dog := Dog{
    Animal: Animal{Name: "Бобик"},
    Breed:  "Лабрадор",
}

fmt.Println(dog.Speak()) // Бобик говорит! — метод продвигается
```

> **Совет:** если нужно переопределить продвинутый метод, просто объявите метод с тем же именем на внешнем типе.
