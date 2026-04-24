---
title: "Площади фигур"
description: "Реализует интерфейс Shape для Circle и Rectangle, считает суммарную площадь"
order: 2
difficulty: medium
---

# Площади фигур

Реализуйте интерфейс `Shape` для двух типов — `Circle` и `Rectangle`, а затем функцию `TotalArea`.

Интерфейс `Shape` уже объявлен в шаблоне. Реализуйте:

1. `Circle` с полем `Radius float64` и методом `Area() float64` → `π * r²`
2. `Rectangle` с полями `Width, Height float64` и методом `Area() float64` → `w * h`
3. `TotalArea(shapes []Shape) float64` — сумма площадей всех фигур

Используйте `math.Pi` для числа π.

## Пример

| Фигура               | `Area()`         |
|----------------------|------------------|
| `Circle{Radius: 5}`  | `≈ 78.54`        |
| `Rectangle{3, 4}`    | `12.0`           |
| `Circle{Radius: 1}`  | `≈ 3.14`         |

`TotalArea([Circle{5}, Rectangle{3,4}])` → `≈ 90.54`
