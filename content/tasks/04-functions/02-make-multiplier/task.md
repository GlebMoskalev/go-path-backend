---
title: "Фабрика множителей"
description: "Возвращает замыкание, умножающее входное число на заданный factor"
order: 2
difficulty: easy
---

# Фабрика множителей

Напишите функцию `MakeMultiplier`, которая принимает число `factor int` и возвращает **функцию**, умножающую своё аргумент на `factor`.

Это классическое замыкание (closure): возвращаемая функция «запоминает» значение `factor` из внешней функции.

## Пример

```go
double := MakeMultiplier(2)
triple := MakeMultiplier(3)

double(5)  // 10
triple(5)  // 15
double(0)  // 0
```

| Вызов              | Выход |
|--------------------|-------|
| `MakeMultiplier(2)(5)`  | `10`  |
| `MakeMultiplier(3)(5)`  | `15`  |
| `MakeMultiplier(0)(99)` | `0`   |
| `MakeMultiplier(2)(0)`  | `0`   |
