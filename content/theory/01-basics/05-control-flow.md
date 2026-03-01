---
title: "Управляющие конструкции"
description: "if/else, switch, for — циклы и условия в Go"
order: 5
---

# Управляющие конструкции

Go имеет минимальный набор управляющих конструкций: `if`, `switch` и единственный цикл `for`.

## if / else

```go
x := 10

if x > 0 {
    fmt.Println("положительное")
} else if x < 0 {
    fmt.Println("отрицательное")
} else {
    fmt.Println("ноль")
}
```

### if с инициализацией

Переменная, объявленная в `if`, видна только внутри блока:

```go
if err := doSomething(); err != nil {
    fmt.Println("ошибка:", err)
}
// err здесь недоступна
```

## switch

```go
day := "Monday"

switch day {
case "Monday":
    fmt.Println("Понедельник")
case "Friday":
    fmt.Println("Пятница!")
default:
    fmt.Println("Обычный день")
}
```

### switch без выражения

Работает как цепочка `if-else`:

```go
score := 85

switch {
case score >= 90:
    fmt.Println("Отлично")
case score >= 70:
    fmt.Println("Хорошо")
default:
    fmt.Println("Нужно подтянуть")
}
```

### fallthrough

В Go `case` не проваливается автоматически. Используйте `fallthrough` явно:

```go
switch 3 {
case 3:
    fmt.Println("три")
    fallthrough
case 4:
    fmt.Println("четыре (тоже выполнится)")
}
```

## Цикл for

В Go только один цикл — `for`. Он заменяет `while` и `do-while`.

### Классический

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

### Как while

```go
n := 1
for n < 100 {
    n *= 2
}
```

### Бесконечный

```go
for {
    // выход через break или return
    break
}
```

### range — перебор коллекций

```go
nums := []int{10, 20, 30}

for i, v := range nums {
    fmt.Printf("индекс %d → значение %d\n", i, v)
}
```

Если индекс не нужен:

```go
for _, v := range nums {
    fmt.Println(v)
}
```

## break и continue

```go
for i := 0; i < 10; i++ {
    if i == 3 {
        continue  // пропустить итерацию
    }
    if i == 7 {
        break     // выйти из цикла
    }
    fmt.Println(i)
}
```

> **Совет:** используйте метки (labels) для выхода из вложенных циклов:
> ```go
> outer:
> for i := 0; i < 3; i++ {
>     for j := 0; j < 3; j++ {
>         if i == 1 && j == 1 {
>             break outer
>         }
>     }
> }
> ```
