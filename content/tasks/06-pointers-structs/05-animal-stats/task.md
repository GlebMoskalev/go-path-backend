---
title: "Встраивание структур"
description: "Реализует Animal с встроенной Stats, демонстрируя продвижение методов"
order: 5
difficulty: hard
---

# Встраивание структур

В шаблоне уже определена структура `Stats` с методом `Summary() string`. Ваша задача — реализовать структуру `Animal` с **встроенной** `Stats` и метод `Describe() string`.

За счёт встраивания `Animal` автоматически получает метод `Summary()` без его явного объявления — это называется «продвижение методов» (method promotion).

Реализуйте:
- Структуру `Animal` со встроенной `Stats` и полем `Name string`
- Метод `Describe() string` — возвращает `"<Name>: <Summary()>"`

## Пример

```go
a := Animal{
    Name: "Лиса",
    Stats: Stats{Speed: 50, Weight: 6},
}
a.Summary()   // "speed=50, weight=6"  — через продвижение!
a.Describe()  // "Лиса: speed=50, weight=6"
```

| Name     | Speed | Weight | `Describe()`                    |
|----------|-------|--------|---------------------------------|
| `"Лиса"` | `50`  | `6`    | `"Лиса: speed=50, weight=6"`    |
| `"Медведь"` | `30` | `200` | `"Медведь: speed=30, weight=200"` |
