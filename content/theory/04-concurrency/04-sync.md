---
title: "Пакет sync"
description: "Mutex, RWMutex, Once, Map — примитивы синхронизации"
order: 4
---

# Пакет sync

Когда каналы не подходят — используйте примитивы синхронизации из пакета `sync`.

## Гонка данных (data race)

```go
// ПРОБЛЕМА: несколько горутин меняют одну переменную
counter := 0
for i := 0; i < 1000; i++ {
    go func() {
        counter++ // гонка данных!
    }()
}
```

Обнаружение: `go run -race main.go`

## sync.Mutex

Мьютекс обеспечивает эксклюзивный доступ к ресурсу:

```go
type SafeCounter struct {
    mu sync.Mutex
    value int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

```go
func main() {
    counter := &SafeCounter{}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Inc()
        }()
    }

    wg.Wait()
    fmt.Println(counter.Value()) // ровно 1000
}
```

## sync.RWMutex

Позволяет множественные одновременные чтения, но эксклюзивную запись:

```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()         // блокировка на чтение
    defer c.mu.RUnlock()
    val, ok := c.data[key]
    return val, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()          // блокировка на запись
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### Когда использовать RWMutex

- Много чтений, мало записей → `RWMutex`
- Примерно равное количество → обычный `Mutex`

## sync.Once

Гарантирует выполнение функции ровно один раз:

```go
var (
    instance *Database
    once     sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        instance = connectToDatabase()
        fmt.Println("подключение создано")
    })
    return instance
}
```

Как бы много горутин ни вызвали `GetDB()`, подключение создаётся один раз.

## sync.Map

Потокобезопасная карта, не требующая внешней блокировки:

```go
var m sync.Map

// Запись
m.Store("key1", "value1")
m.Store("key2", 42)

// Чтение
val, ok := m.Load("key1")
if ok {
    fmt.Println(val) // value1
}

// Удаление
m.Delete("key1")

// Чтение или запись
actual, loaded := m.LoadOrStore("key3", "default")
// loaded=false → записано "default"
// loaded=true  → вернуло существующее значение

// Перебор
m.Range(func(key, value any) bool {
    fmt.Printf("%v: %v\n", key, value)
    return true // false — остановить перебор
})
```

### Когда использовать sync.Map

- Ключи стабильны (мало записей, много чтений)
- Разные горутины работают с разными ключами
- Для обычных случаев `map` + `RWMutex` эффективнее

## sync.Pool

Пул переиспользуемых объектов для снижения нагрузки на GC:

```go
var bufPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer)
    },
}

func process(data string) string {
    buf := bufPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufPool.Put(buf)
    }()

    buf.WriteString("processed: ")
    buf.WriteString(data)
    return buf.String()
}
```

## atomic — атомарные операции

Для простых счётчиков мьютексы избыточны:

```go
import "sync/atomic"

var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}

func value() int64 {
    return atomic.LoadInt64(&counter)
}
```

Go 1.19+ предоставляет типизированные атомики:

```go
var counter atomic.Int64

counter.Add(1)
fmt.Println(counter.Load()) // 1
```

> **Совет:** выбирайте подходящий инструмент:
> - Координация → каналы
> - Защита данных → `Mutex` / `RWMutex`
> - Простые счётчики → `atomic`
> - Одноразовая инициализация → `sync.Once`
