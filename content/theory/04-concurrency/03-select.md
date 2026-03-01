---
title: "Select"
description: "Мультиплексирование каналов, таймауты, неблокирующие операции"
order: 3
---

# Select

`select` позволяет горутине ждать на нескольких каналах одновременно.

## Базовый синтаксис

```go
select {
case msg := <-ch1:
    fmt.Println("из ch1:", msg)
case msg := <-ch2:
    fmt.Println("из ch2:", msg)
}
```

`select` блокируется, пока один из каналов не станет готов. Если готовы несколько — выбирается случайный.

## Пример с двумя источниками

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "один"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "два"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg := <-ch1:
            fmt.Println("Получено:", msg)
        case msg := <-ch2:
            fmt.Println("Получено:", msg)
        }
    }
}
```

## Таймаут

```go
select {
case result := <-longOperation():
    fmt.Println("результат:", result)
case <-time.After(3 * time.Second):
    fmt.Println("таймаут!")
}
```

## default — неблокирующий select

```go
select {
case msg := <-ch:
    fmt.Println("получено:", msg)
default:
    fmt.Println("канал пуст, продолжаем")
}
```

### Неблокирующая отправка

```go
select {
case ch <- value:
    fmt.Println("отправлено")
default:
    fmt.Println("канал полон, пропускаем")
}
```

## Бесконечный цикл с select

```go
func worker(done <-chan struct{}, tasks <-chan int) {
    for {
        select {
        case <-done:
            fmt.Println("завершение")
            return
        case task := <-tasks:
            fmt.Println("обработка:", task)
        }
    }
}
```

## Паттерн: done-канал

```go
func main() {
    done := make(chan struct{})
    go worker(done)

    time.Sleep(3 * time.Second)
    close(done) // сигнал завершения всем горутинам
}

func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            return
        default:
            // работа...
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

## Паттерн: тикер

```go
func main() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    done := make(chan struct{})

    go func() {
        time.Sleep(5 * time.Second)
        close(done)
    }()

    for {
        select {
        case <-done:
            fmt.Println("стоп")
            return
        case t := <-ticker.C:
            fmt.Println("тик:", t.Format("15:04:05"))
        }
    }
}
```

## context.Context — продвинутая отмена

```go
func main() {
    ctx, cancel := context.WithTimeout(
        context.Background(),
        3*time.Second,
    )
    defer cancel()

    select {
    case result := <-doWork(ctx):
        fmt.Println("результат:", result)
    case <-ctx.Done():
        fmt.Println("отмена:", ctx.Err())
    }
}

func doWork(ctx context.Context) <-chan string {
    ch := make(chan string)
    go func() {
        // долгая операция...
        time.Sleep(5 * time.Second)
        ch <- "готово"
    }()
    return ch
}
```

> **Совет:** используйте `context.Context` вместо `done`-каналов — он поддерживает таймауты, дедлайны и передачу значений.
