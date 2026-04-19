---
title: "Горутины"
description: "go keyword, планировщик Go (M:N threading, GOMAXPROCS), goroutine leak"
order: 1
---

# Горутины

Горутина — лёгкий поток выполнения, управляемый планировщиком Go. Это фундамент конкурентности в Go.

## Запуск горутины

```go
go функция(аргументы)
```

```go
func main() {
    go fmt.Println("из горутины")
    fmt.Println("из main")
    time.Sleep(time.Millisecond) // дать горутине выполниться
}
// из main
// из горутины  (порядок не гарантирован)
```

Анонимная функция:

```go
go func() {
    fmt.Println("анонимная горутина")
}()  // обязательно вызвать немедленно через ()
```

С замыканием:

```go
msg := "hello"
go func() {
    fmt.Println(msg)  // захватывает msg
}()
```

---

## Планировщик Go: M:N threading

Go использует модель M:N threading: M горутин выполняются на N OS-потоках. Планировщик Go (runtime scheduler) распределяет горутины по потокам.

```
Горутины (M):  G1  G2  G3  G4  G5  G6
                ↕   ↕   ↕
OS-потоки (N): T1  T2  T3
                ↕   ↕   ↕
CPU-ядра (P):  P1  P2  P3
```

**Преимущества:**
- Горутины стартуют с маленьким стеком (~2-8 KB), который растёт по мере необходимости
- Создание тысяч горутин — нормально (в отличие от OS-потоков, где каждый ~1-8 MB стека)
- Кооперативное + преемптивное планирование (Go 1.14+)

### GOMAXPROCS

Количество OS-потоков для выполнения горутин. По умолчанию = количество CPU-ядер:

```go
import "runtime"

fmt.Println(runtime.GOMAXPROCS(0))  // текущее значение
runtime.GOMAXPROCS(4)               // установить вручную

fmt.Println(runtime.NumCPU())       // количество физических ядер
fmt.Println(runtime.NumGoroutine()) // текущее количество горутин
```

---

## Goroutine leak — утечка горутин

Горутина «течёт» когда она запущена, но никогда не завершится. Накопившиеся горутины потребляют память и могут исчерпать ресурсы.

### Частые причины утечек

**Чтение из канала, который никогда не закрывается:**

```go
// УТЕЧКА: горутина ждёт вечно
func leaky(ch <-chan int) {
    go func() {
        for v := range ch {  // если ch никогда не закроют...
            fmt.Println(v)
        }
    }()
}
```

**Запись в канал без читателя:**

```go
// УТЕЧКА: горутина заблокирована навсегда
func leaky2() {
    ch := make(chan int)
    go func() {
        ch <- 42  // никто не читает
    }()
}
```

**HTTP-запрос без контекста:**

```go
// УТЕЧКА: горутина выполняет запрос даже после отмены операции
func startRequest() {
    go func() {
        resp, _ := http.Get("https://slow-api.example.com/data")
        // обработка...
        _ = resp
    }()
}
```

### Как обнаружить утечку

**runtime.NumGoroutine() в тестах:**

```go
func TestNoLeak(t *testing.T) {
    before := runtime.NumGoroutine()
    
    // ... запускаем код ...
    doSomething()
    
    time.Sleep(10 * time.Millisecond) // дать время завершиться
    after := runtime.NumGoroutine()
    
    if after > before {
        t.Errorf("горутин утекло: было %d, стало %d", before, after)
    }
}
```

**Пакет goleak (golangci-lint):**

```go
import "go.uber.org/goleak"

func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}
```

### Как предотвратить утечки

**1. Done-канал для отмены:**

```go
func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            return  // корректное завершение
        default:
            doWork()
        }
    }
}

done := make(chan struct{})
go worker(done)
// ...
close(done)  // сигнал завершения
```

**2. context.Context:**

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            doWork()
        }
    }
}

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
go worker(ctx)
```

> 💡 context.Context подробно разберём в главе 12 «Продвинутые темы».

---

## sync.WaitGroup — ожидание горутин

Самый распространённый способ дождаться завершения группы горутин:

```go
import "sync"

func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1)  // увеличить счётчик перед запуском горутины
        go func(i int) {
            defer wg.Done()  // уменьшить счётчик при завершении
            fmt.Println("горутина", i)
        }(i)
    }
    
    wg.Wait()  // ждать пока счётчик не станет 0
    fmt.Println("все горутины завершились")
}
```

**Правило**: всегда вызывай `wg.Add(1)` **до** запуска горутины, не внутри неё.

---

## Горутины и замыкания: ловушка

```go
// НЕВЕРНО: все горутины видят последнее значение i
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // вероятно напечатает 5 5 5 5 5
    }()
}

// ВЕРНО: передаём i как параметр
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n)  // 0 1 2 3 4 (в произвольном порядке)
    }(i)
}
```

> Подробнее про ловушку замыканий в цикле — в главе 4 «Функции».

---

## Итог

- `go f()` запускает f в новой горутине; управление немедленно возвращается
- Планировщик Go: M горутин на N OS-потоках (M:N threading)
- Горутины легковесны — создавай тысячи, не OS-потоки
- GOMAXPROCS = количество параллельных потоков (по умолчанию = кол-во CPU)
- Горутина течёт, когда никогда не завершается: используй done-канал или context
- `sync.WaitGroup` — для ожидания группы горутин; `Add` перед запуском, `Done` в defer
