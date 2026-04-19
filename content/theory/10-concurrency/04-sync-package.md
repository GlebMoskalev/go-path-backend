---
title: "Пакет sync"
description: "WaitGroup, Mutex, RWMutex, Once, atomic"
order: 4
---

# Пакет sync

Каналы — не единственный способ синхронизации в Go. Пакет `sync` предоставляет примитивы для случаев, когда горутинам нужен общий доступ к памяти.

## sync.WaitGroup

`WaitGroup` ждёт завершения группы горутин. Мы уже рассматривали его в главе о горутинах, здесь углубимся в детали.

```go
import "sync"

var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)           // перед запуском
    go func(n int) {
        defer wg.Done() // при завершении
        fmt.Println(n)
    }(i)
}

wg.Wait() // блокируется пока счётчик не станет 0
```

### Типичные ошибки WaitGroup

**Add внутри горутины (гонка с Wait):**

```go
// НЕВЕРНО: Add может вызваться после Wait
for i := 0; i < 5; i++ {
    go func(n int) {
        wg.Add(1)       // ОШИБКА: слишком поздно
        defer wg.Done()
        fmt.Println(n)
    }(i)
}
wg.Wait()
```

**Передача по значению (копия не синхронизируется):**

```go
// НЕВЕРНО: wg скопирован
func doWork(wg sync.WaitGroup) { // должно быть *sync.WaitGroup
    defer wg.Done()
    // ...
}
```

Всегда передавай WaitGroup через указатель.

---

## sync.Mutex — взаимное исключение

`Mutex` защищает данные от одновременного доступа нескольких горутин:

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    c := &Counter{}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Increment()
        }()
    }

    wg.Wait()
    fmt.Println(c.Value()) // гарантированно 1000
}
```

### Правила Mutex

- `Lock()` блокируется если мьютекс уже захвачен другой горутиной
- `Unlock()` в `defer` — обязательная практика, защита от паники
- Никогда не копируй `sync.Mutex` (и любые типы из пакета sync) после первого использования

```go
// НЕВЕРНО: копирование мьютекса
m1 := sync.Mutex{}
m2 := m1  // m2 — независимая копия внутреннего состояния, не синхронизирована
```

Инструмент `go vet` обнаружит такие копирования.

### TryLock (Go 1.18+)

```go
if mu.TryLock() {
    defer mu.Unlock()
    // успешно захватили
} else {
    // мьютекс занят, делаем что-то другое
}
```

---

## sync.RWMutex — читатели/писатели

Когда чтений много, а записей мало, `RWMutex` эффективнее обычного `Mutex`: несколько горутин могут читать одновременно.

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()           // эксклюзивная блокировка для записи
    defer c.mu.Unlock()
    c.items[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()          // разделяемая блокировка для чтения
    defer c.mu.RUnlock()
    v, ok := c.items[key]
    return v, ok
}
```

**Когда использовать RWMutex:**
- Много операций чтения, редкие записи
- Критическая секция при чтении не изменяет данные

Если записи и чтения примерно поровну — обычный Mutex может быть быстрее (меньше накладных расходов).

---

## sync.Once — однократное выполнение

`Once` гарантирует что функция выполнится ровно один раз, даже при конкурентных вызовах. Классический сценарий — ленивая инициализация:

```go
type Singleton struct {
    data string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        fmt.Println("инициализация (один раз)")
        instance = &Singleton{data: "данные"}
    })
    return instance
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            s := GetInstance()
            fmt.Println(s.data)
        }()
    }
    wg.Wait()
}
// инициализация (один раз)
// данные
// данные
// данные
// данные
// данные
```

### Нюанс: паника внутри Do

Если функция в `Do` запаникует — `Once` всё равно считается выполненным. Повторный вызов `Do` не выполнит функцию снова:

```go
var once sync.Once
var initialized bool

once.Do(func() {
    panic("ошибка инициализации")
})

once.Do(func() {
    initialized = true  // НИКОГДА не выполнится
})
```

---

## sync.Map — конкурентная карта

Стандартный `map` не безопасен для конкурентного доступа. `sync.Map` — безопасный вариант с особым API:

```go
var m sync.Map

// Запись:
m.Store("key", "value")

// Чтение:
v, ok := m.Load("key")
if ok {
    fmt.Println(v.(string))
}

// Загрузить или сохранить:
actual, loaded := m.LoadOrStore("key", "default")
// loaded=true если ключ уже был, false если сохранили default

// Удалить:
m.Delete("key")

// Обход:
m.Range(func(k, v any) bool {
    fmt.Println(k, v)
    return true  // false — остановить обход
})
```

**Когда использовать sync.Map:**
- Ключи пишутся один раз, читаются многократно
- Горутины работают с разными ключами (нет конкуренции за один ключ)

В остальных случаях — `map` + `sync.Mutex` проще и быстрее.

---

## sync/atomic — атомарные операции

Пакет `sync/atomic` предоставляет атомарные операции для примитивных типов без блокировок:

```go
import "sync/atomic"

var counter int64

// Атомарное увеличение:
atomic.AddInt64(&counter, 1)

// Атомарное чтение:
v := atomic.LoadInt64(&counter)

// Атомарная запись:
atomic.StoreInt64(&counter, 42)

// Сравнить и заменить (Compare-and-Swap):
swapped := atomic.CompareAndSwapInt64(&counter, 42, 100)
// если counter == 42, устанавливает 100 и возвращает true
```

### atomic.Value — атомарное хранилище произвольных типов

```go
var config atomic.Value

// Записать:
cfg := Config{MaxConns: 10}
config.Store(cfg)

// Прочитать:
v := config.Load()
if v != nil {
    cfg := v.(Config)
    fmt.Println(cfg.MaxConns)
}
```

Используется для атомарной горячей замены конфигурации.

**Atomic vs Mutex:**
- Atomic: только для примитивов, максимальная скорость, нет ожидания
- Mutex: для сложных структур данных, явная блокировка секции

---

## Типичные ошибки

### Дедлок

Дедлок возникает когда горутины циклически ждут друг друга:

```go
// ДЕДЛОК: Lock два раза из одной горутины
mu.Lock()
mu.Lock()  // блокируется навсегда: sync.Mutex не рекурентный
```

Go runtime обнаруживает полные дедлоки и завершает программу с `fatal error: all goroutines are asleep - deadlock!`.

### Слишком широкая блокировка

```go
// НЕВЕРНО: держим мьютекс во время долгой операции
mu.Lock()
result, err := http.Get(url)  // HTTP-запрос под мьютексом!
mu.Unlock()

// ВЕРНО: минимальная критическая секция
resp, err := http.Get(url)   // без мьютекса
if err == nil {
    mu.Lock()
    cache[url] = resp
    mu.Unlock()
}
```

---

## Итог

| Примитив | Назначение |
|----------|-----------|
| `WaitGroup` | Ждать группу горутин; Add перед go, Done в defer |
| `Mutex` | Защита данных от конкурентного доступа |
| `RWMutex` | Много читателей, мало писателей |
| `Once` | Однократная инициализация |
| `sync.Map` | Конкурентная карта для специфических сценариев |
| `atomic` | Атомарные операции с числами без блокировок |

- Никогда не копируй типы из пакета `sync` после первого использования
- `defer mu.Unlock()` — всегда, защита от паники
- Минимизируй время под блокировкой
