---
title: "Race Detector"
description: "-race флаг, гонки данных, sync/atomic для счётчиков"
order: 6
---

# Race Detector

Гонка данных (data race) — одна из самых коварных ошибок в конкурентных программах. Go поставляется со встроенным детектором гонок.

## Что такое гонка данных

Гонка данных возникает когда две горутины обращаются к одной переменной одновременно, и хотя бы одна из них пишет:

```go
// ГОНКА ДАННЫХ
var counter int

func main() {
    for i := 0; i < 1000; i++ {
        go func() {
            counter++  // одновременное чтение и запись из разных горутин
        }()
    }
    time.Sleep(time.Second)
    fmt.Println(counter)  // непредсказуемый результат: 947, 1000, 823...
}
```

Результат недетерминирован — зависит от планировщика, процессора и кэшей.

---

## Флаг -race

Запусти программу с флагом `-race`, и Go runtime обнаружит гонки:

```bash
go run -race main.go
go test -race ./...
go build -race -o myapp .
```

Пример вывода детектора:

```
==================
WARNING: DATA RACE
Write at 0x00c0000b4010 by goroutine 7:
  main.main.func1()
      /tmp/main.go:8 +0x2c

Previous write at 0x00c0000b4010 by goroutine 6:
  main.main.func1()
      /tmp/main.go:8 +0x2c

Goroutine 7 (running) created at:
  main.main()
      /tmp/main.go:7 +0x3a

Goroutine 6 (running) created at:
  main.main()
      /tmp/main.go:7 +0x3a
==================
```

Детектор показывает:
- Тип операции (Read/Write)
- Адрес памяти
- Стек вызовов каждой горутины
- Где горутины были созданы

### Накладные расходы -race

- Замедление: в 2–20 раз
- Память: в 5–10 раз больше

Использую `-race` при разработке и в CI, но не в production.

---

## Исправление гонок

### Вариант 1: sync/atomic

Для простых счётчиков атомарные операции быстрее мьютекса:

```go
import "sync/atomic"

var counter int64

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1)  // атомарно, нет гонки
        }()
    }
    wg.Wait()
    fmt.Println(atomic.LoadInt64(&counter))  // всегда 1000
}
```

### Вариант 2: sync.Mutex

Для сложных операций или нескольких взаимосвязанных переменных:

```go
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Get() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

### Вариант 3: каналы

Инкапсулировать состояние в одной горутине-«владельце»:

```go
type Counter struct {
    inc chan struct{}
    get chan int
}

func NewCounter() *Counter {
    c := &Counter{
        inc: make(chan struct{}),
        get: make(chan int),
    }
    go func() {
        value := 0
        for {
            select {
            case <-c.inc:
                value++
            case c.get <- value:
            }
        }
    }()
    return c
}

func (c *Counter) Inc() { c.inc <- struct{}{} }
func (c *Counter) Get() int { return <-c.get }
```

Единственная горутина владеет `value` — гонки невозможны по определению.

---

## sync/atomic подробнее

```go
import "sync/atomic"

// Целые числа:
var n int64
atomic.AddInt64(&n, 1)
atomic.AddInt64(&n, -1)   // декремент
v := atomic.LoadInt64(&n)
atomic.StoreInt64(&n, 42)

// Compare-And-Swap:
old, new_ := int64(42), int64(100)
swapped := atomic.CompareAndSwapInt64(&n, old, new_)
// если n == 42, устанавливает 100 и возвращает true

// atomic.Value для произвольных типов:
var cfg atomic.Value
cfg.Store(Config{MaxConns: 10})   // запись
c := cfg.Load().(Config)           // чтение (type assertion)
```

### Паттерн: горячая замена конфигурации

```go
type Config struct {
    MaxConns int
    Timeout  time.Duration
}

var currentConfig atomic.Value

func init() {
    currentConfig.Store(Config{MaxConns: 10, Timeout: 5 * time.Second})
}

// В фоновой горутине:
func reloadConfig() {
    newCfg := loadFromFile()
    currentConfig.Store(newCfg)  // атомарная замена
}

// В обработчике:
func handleRequest() {
    cfg := currentConfig.Load().(Config)
    // использует актуальную конфигурацию без блокировок
}
```

---

## Частые причины гонок

### Замыкание над переменной цикла

```go
// ГОНКА: все горутины читают одну и ту же переменную i
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(i)  // гонка с обновлением i в цикле
    }()
}

// ВЕРНО: передать значение через параметр
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println(n)
    }(i)
}
```

### Возврат указателя на разделяемые данные

```go
type Cache struct {
    mu   sync.Mutex
    data map[string][]int
}

// НЕВЕРНО: возвращаем указатель на внутренний слайс
func (c *Cache) GetRaw(key string) []int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.data[key]  // вызывающий может изменить слайс без блокировки!
}

// ВЕРНО: возвращаем копию
func (c *Cache) Get(key string) []int {
    c.mu.Lock()
    defer c.mu.Unlock()
    src := c.data[key]
    dst := make([]int, len(src))
    copy(dst, src)
    return dst
}
```

### map без синхронизации

```go
// ГОНКА: конкурентный доступ к map
m := make(map[string]int)

go func() { m["a"] = 1 }()
go func() { m["b"] = 2 }()

// Go runtime обнаружит это и завершит программу:
// fatal error: concurrent map writes
```

Начиная с Go 1.6 runtime детектирует конкурентные записи в map и падает с паникой даже без `-race`. Используй `sync.Mutex` или `sync.Map`.

---

## Итог

- Гонка данных: одновременный доступ к переменной из нескольких горутин, хотя бы одна пишет
- `-race` флаг: детектирует гонки в runtime; используй при разработке и в CI
- `sync/atomic`: для простых числовых счётчиков без блокировок
- `sync.Mutex`: для сложной логики и взаимосвязанных данных
- Конкурентный доступ к `map` без синхронизации — фатальная ошибка
- Не возвращай указатели на защищённые данные — возвращай копии
