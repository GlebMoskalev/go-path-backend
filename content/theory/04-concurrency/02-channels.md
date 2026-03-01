---
title: "Каналы"
description: "Каналы: буферизованные и небуферизованные, направление, закрытие"
order: 2
---

# Каналы

Каналы — основной механизм коммуникации между горутинами.

> **Принцип Go:** «Не общайтесь через разделяемую память — делитесь памятью через общение».

## Создание

```go
ch := make(chan int)       // небуферизованный
ch := make(chan int, 10)   // буферизованный на 10 элементов
```

## Отправка и приём

```go
ch <- 42      // отправить значение в канал
value := <-ch // принять значение из канала
```

## Небуферизованный канал

Отправитель блокируется, пока получатель не примет значение:

```go
func main() {
    ch := make(chan string)

    go func() {
        ch <- "привет" // блокируется, пока main не прочитает
    }()

    msg := <-ch // блокируется, пока горутина не отправит
    fmt.Println(msg)
}
```

## Буферизованный канал

Отправитель блокируется, только когда буфер полон:

```go
ch := make(chan int, 3)

ch <- 1  // не блокируется
ch <- 2  // не блокируется
ch <- 3  // не блокируется
// ch <- 4 // заблокируется, буфер полон

fmt.Println(<-ch) // 1
fmt.Println(<-ch) // 2
```

## Закрытие канала

```go
close(ch)
```

После закрытия:
- Чтение возвращает оставшиеся значения, затем нулевое значение
- Запись вызывает **panic**

### Проверка закрытия

```go
value, ok := <-ch
if !ok {
    fmt.Println("канал закрыт")
}
```

## range по каналу

```go
func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch) // ВАЖНО: закрыть канал
}

func main() {
    ch := make(chan int)
    go producer(ch)

    for v := range ch {
        fmt.Println(v) // 0, 1, 2, 3, 4
    }
}
```

## Направление каналов

Ограничивайте использование канала в параметрах функций:

```go
func send(ch chan<- int) {   // только отправка
    ch <- 42
}

func receive(ch <-chan int) { // только приём
    fmt.Println(<-ch)
}
```

## Паттерн: генератор

```go
func fibonacci(n int) <-chan int {
    ch := make(chan int)
    go func() {
        a, b := 0, 1
        for i := 0; i < n; i++ {
            ch <- a
            a, b = b, a+b
        }
        close(ch)
    }()
    return ch
}

for v := range fibonacci(10) {
    fmt.Println(v)
}
```

## Паттерн: fan-out / fan-in

```go
// Fan-out: несколько горутин читают из одного канала
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Запускаем 3 воркера
    for w := 0; w < 3; w++ {
        go worker(w, jobs, results)
    }

    // Отправляем задачи
    for j := 0; j < 9; j++ {
        jobs <- j
    }
    close(jobs)

    // Собираем результаты
    for r := 0; r < 9; r++ {
        fmt.Println(<-results)
    }
}
```

> **Совет:** всегда думайте о том, кто закроет канал. Обычно это делает отправитель.
