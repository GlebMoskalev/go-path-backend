---
title: "Select"
description: "select statement, default case, таймаут с time.After, fan-in паттерн"
order: 3
---

# Select

`select` позволяет горутине ждать сразу нескольких операций с каналами. Это аналог `switch`, но для каналов.

## Синтаксис

```go
select {
case v := <-ch1:
    fmt.Println("получено из ch1:", v)
case ch2 <- value:
    fmt.Println("отправлено в ch2")
case v, ok := <-ch3:
    if !ok {
        fmt.Println("ch3 закрыт")
    }
}
```

`select` блокируется пока ни один из каналов не готов. Когда несколько каналов готовы одновременно — выбирается **случайный** (псевдослучайный) case.

---

## Базовый пример

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
            fmt.Println("ch1:", msg)
        case msg := <-ch2:
            fmt.Println("ch2:", msg)
        }
    }
}
// ch1: один
// ch2: два
```

---

## Default case — неблокирующий select

`default` выполняется когда ни один канал не готов:

```go
ch := make(chan int, 1)

select {
case v := <-ch:
    fmt.Println("получено:", v)
default:
    fmt.Println("канал пуст, идём дальше")
}
```

Неблокирующая отправка:

```go
select {
case ch <- value:
    fmt.Println("отправлено")
default:
    fmt.Println("буфер полон, пропускаем")
}
```

### Паттерн: проверка отмены без блокировки

```go
func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            fmt.Println("завершаем работу")
            return
        default:
            // продолжаем работу
            doWork()
        }
    }
}
```

---

## Таймаут с time.After

`time.After(d)` возвращает канал, в который придёт время через `d`:

```go
func fetchWithTimeout(ch <-chan string) (string, error) {
    select {
    case result := <-ch:
        return result, nil
    case <-time.After(2 * time.Second):
        return "", fmt.Errorf("таймаут: нет ответа за 2 секунды")
    }
}
```

Пример с HTTP-запросом:

```go
func main() {
    resultCh := make(chan string, 1)

    go func() {
        // имитация долгого запроса
        time.Sleep(3 * time.Second)
        resultCh <- "данные получены"
    }()

    select {
    case result := <-resultCh:
        fmt.Println(result)
    case <-time.After(1 * time.Second):
        fmt.Println("ошибка: таймаут")
    }
}
// ошибка: таймаут
```

### time.NewTimer vs time.After

`time.After` создаёт таймер, который не освобождается до истечения времени (риск утечки памяти в циклах). В цикле используй `time.NewTimer`:

```go
timer := time.NewTimer(2 * time.Second)
defer timer.Stop()

select {
case result := <-ch:
    timer.Stop()
    use(result)
case <-timer.C:
    fmt.Println("таймаут")
}
```

---

## Ticker — периодические события

```go
func main() {
    ticker := time.NewTicker(500 * time.Millisecond)
    done := make(chan struct{})

    go func() {
        time.Sleep(2 * time.Second)
        close(done)
    }()

    for {
        select {
        case t := <-ticker.C:
            fmt.Println("тик:", t.Format("15:04:05.000"))
        case <-done:
            ticker.Stop()
            fmt.Println("стоп")
            return
        }
    }
}
// тик: 14:32:00.500
// тик: 14:32:01.000
// тик: 14:32:01.500
// тик: 14:32:02.000
// стоп
```

---

## Fan-in паттерн

Fan-in — паттерн для слияния нескольких каналов в один. Полезен когда несколько горутин производят данные параллельно, а потребитель хочет читать из одного места.

Без fan-in пришлось бы читать каналы по очереди и терять параллелизм. С fan-in — читаем из merged и получаем данные сразу со всех источников, в том порядке в котором они готовы.

```go
func fanIn(ch1, ch2 <-chan string) <-chan string {
    out := make(chan string)

    go func() {
        defer close(out)
        for {
            select {
            case v, ok := <-ch1:
                if !ok {
                    ch1 = nil  // канал закрыт — отключаем его case
                } else {
                    out <- v
                }
            case v, ok := <-ch2:
                if !ok {
                    ch2 = nil  // канал закрыт — отключаем его case
                } else {
                    out <- v
                }
            }
            if ch1 == nil && ch2 == nil {
                return  // оба закрыты — завершаем горутину
            }
        }
    }()

    return out
}
```

Ключевой трюк: когда канал закрылся, мы присваиваем ему `nil`. Case с nil-каналом в select **никогда не выбирается** — операция на nil-канале блокируется навсегда, и select её игнорирует. Это позволяет корректно дождаться закрытия обоих каналов.

Пример с двумя источниками, которые отправляют данные с паузой:

```go
func source(name string, n int) <-chan string {
    ch := make(chan string)
    go func() {
        defer close(ch)
        for i := 0; i < n; i++ {
            ch <- fmt.Sprintf("%s-%d", name, i)
            time.Sleep(100 * time.Millisecond)
        }
    }()
    return ch
}

func main() {
    merged := fanIn(source("A", 3), source("B", 3))
    for v := range merged {
        fmt.Println(v)
    }
}
```

Возможный вывод:
```
A-0
B-0
A-1
B-1
A-2
B-2
```

Порядок зависит от планировщика и нагрузки. В данном случае A и B стартуют одновременно и делают паузы одинаковой длины, поэтому чередуются почти ровно. На другой машине или под нагрузкой порядок может отличаться — это нормально, ведь select выбирает готовый case случайно если их несколько.

---

## Nil-канал как механизм управления

`select` с nil-каналом всегда блокируется в этом case — полезно для условного включения/отключения:

```go
func merge(a, b <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for a != nil || b != nil {
            select {
            case v, ok := <-a:
                if !ok {
                    a = nil
                    continue
                }
                out <- v
            case v, ok := <-b:
                if !ok {
                    b = nil
                    continue
                }
                out <- v
            }
        }
    }()
    return out
}
```

---

## Типичные ошибки

### Случайность выбора case

Если оба канала готовы, Go выбирает случайно. Это корректное поведение, но неочевидное:

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)
ch1 <- 1
ch2 <- 2

select {
case v := <-ch1:
    fmt.Println("ch1:", v)
case v := <-ch2:
    fmt.Println("ch2:", v)
}
// может напечатать либо ch1:1, либо ch2:2
```

Если важен порядок — проверяй последовательно через `if`.

### Busy loop и пустой default

Busy loop — это цикл, который крутится без остановки, не делая полезной работы и не отдавая управление. Он нагружает CPU на 100%, даже если реально ждёт какое-то событие.

`select` без `default` ведёт себя правильно: если ни один канал не готов — горутина засыпает и не тратит CPU. Это нормально.

```go
// OK: горутина спит пока ch не даст значение
for {
    select {
    case v := <-ch:
        process(v)
    }
}
```

Пустой `default` ломает это поведение — теперь если канал не готов, select мгновенно проваливается в `default` и цикл тут же начинает следующую итерацию. Горутина не засыпает, а бесконечно проверяет канал снова и снова:

```go
// ПЛОХО: горутина крутится вхолостую и грузит CPU на 100%
for {
    select {
    case v := <-ch:
        process(v)
    default:
        // ничего не делаем, но CPU уже занят
    }
}
```

Добавляй `default` только когда есть реальная работа при пустом канале — например, неблокирующая попытка чтения перед тем как заняться чем-то другим.

---

## Итог

- `select` ждёт готовности одного из каналов; при нескольких готовых выбирает случайный
- `default` делает select неблокирующим
- `time.After(d)` — таймаут в select; `time.NewTimer` предпочтительнее в циклах
- Nil-канал в select никогда не выбирается — используй для условного отключения
- Fan-in: несколько входных каналов → один выходной через select с nil-отключением
