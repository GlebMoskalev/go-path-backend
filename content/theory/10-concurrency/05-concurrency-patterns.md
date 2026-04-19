---
title: "Паттерны конкурентности"
description: "Worker pool, pipeline, fan-out/fan-in, context.Context для отмены"
order: 5
---

# Паттерны конкурентности

Горутины и каналы — строительные блоки. Паттерны — проверенные способы их комбинирования для решения типичных задач.

## Worker Pool — пул воркеров

Ограничивает количество параллельных операций. Полезно когда задач много, но ресурсы (CPU, память, соединения) ограничены.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID int
}

type Result struct {
    JobID  int
    Output int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        // имитация работы
        time.Sleep(100 * time.Millisecond)
        results <- Result{JobID: job.ID, Output: job.ID * 2}
        fmt.Printf("воркер %d обработал задачу %d\n", id, job.ID)
    }
}

func main() {
    const numJobs    = 20
    const numWorkers = 5

    jobs    := make(chan Job, numJobs)
    results := make(chan Result, numJobs)

    var wg sync.WaitGroup
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    // Отправить все задачи
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j}
    }
    close(jobs)  // сигнал: задачи закончились

    // Ждать воркеров и закрыть results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Собрать результаты
    for r := range results {
        _ = r
    }
    fmt.Println("все задачи выполнены")
}
```

---

## Pipeline — конвейер

Цепочка стадий обработки, каждая выполняется в отдельной горутине:

```go
// Стадия 1: генерация чисел
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// Стадия 2: возведение в квадрат
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// Стадия 3: фильтрация (только чётные)
func filterEven(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if n%2 == 0 {
                out <- n
            }
        }
    }()
    return out
}

func main() {
    // Соединяем конвейер:
    nums := generate(1, 2, 3, 4, 5, 6)
    squared := square(nums)
    evens := filterEven(squared)

    for v := range evens {
        fmt.Println(v)
    }
    // 4 16 36
}
```

Каждая стадия читает из входного канала и пишет в выходной. Каналы автоматически синхронизируют скорость стадий.

---

## Fan-out — распределение работы

Один канал-источник, несколько горутин-обработчиков:

```go
func fanOut(in <-chan int, workers int) []<-chan int {
    outs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        out := make(chan int)
        outs[i] = out
        go func(out chan<- int) {
            defer close(out)
            for v := range in {
                out <- v * v
            }
        }(out)
    }
    return outs
}
```

Но обычно fan-out реализуется через worker pool, а fan-in — отдельно.

---

## Отмена через context.Context

`context.Context` — стандартный способ распространения отмены, таймаутов и значений через цепочку горутин.

> 💡 context.Context подробно разберём в главе 12 «Продвинутые темы». Здесь — практический минимум для паттернов конкурентности.

```go
import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int, jobs <-chan int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("воркер %d: отмена (%v)\n", id, ctx.Err())
            return
        case job, ok := <-jobs:
            if !ok {
                return
            }
            // имитация работы
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("воркер %d: задача %d\n", id, job)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    jobs := make(chan int, 10)
    for i := 1; i <= 10; i++ {
        jobs <- i
    }
    close(jobs)

    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(ctx, id, jobs)
        }(i)
    }
    wg.Wait()
}
```

### WithCancel — ручная отмена

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    // через 2 секунды отменяем всё
    time.Sleep(2 * time.Second)
    cancel()
}()

// В воркере:
select {
case <-ctx.Done():
    return  // отмена получена
case result := <-work:
    process(result)
}
```

---

## Паттерн: Done-канал с context

Конвейер с поддержкой отмены:

```go
func generateWithCtx(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-ctx.Done():
                return  // отменили — прекращаем
            }
        }
    }()
    return out
}

func squareWithCtx(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    nums := generateWithCtx(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squared := squareWithCtx(ctx, nums)

    // Берём только первые 3 значения
    for i := 0; i < 3; i++ {
        fmt.Println(<-squared)
    }
    // cancel() освободит все горутины конвейера
}
```

---

## Семафор через буферизованный канал

Ограничение количества одновременных операций:

```go
type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore {
    return make(Semaphore, n)
}

func (s Semaphore) Acquire() { s <- struct{}{} }
func (s Semaphore) Release() { <-s }

func main() {
    sem := NewSemaphore(3)  // максимум 3 параллельно
    var wg sync.WaitGroup

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            sem.Acquire()
            defer sem.Release()

            fmt.Printf("задача %d начата\n", id)
            time.Sleep(500 * time.Millisecond)
            fmt.Printf("задача %d завершена\n", id)
        }(i)
    }
    wg.Wait()
}
```

---

## Итог

| Паттерн | Назначение |
|---------|-----------|
| Worker Pool | Ограничить параллелизм; N воркеров берут задачи из канала |
| Pipeline | Цепочка стадий обработки; каналы синхронизируют скорость |
| Fan-out | Распределить задачи по нескольким горутинам |
| Fan-in | Объединить несколько каналов в один |
| Семафор | Ограничить количество одновременных операций |
| context | Распространить отмену вниз по цепочке горутин |

Основной принцип: горутина должна **всегда иметь способ завершиться** — через закрытый канал, context или done-сигнал.
