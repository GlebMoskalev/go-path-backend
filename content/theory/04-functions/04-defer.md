---
title: "defer"
description: "Порядок выполнения LIFO, defer с panic/recover, cleanup паттерны"
order: 4
---

# defer

`defer` откладывает выполнение функции до момента выхода из текущей функции. Это элегантный механизм для гарантированной очистки ресурсов.

## Базовый синтаксис и порядок выполнения

```go
func main() {
    defer fmt.Println("три")
    defer fmt.Println("два")
    defer fmt.Println("один")
    fmt.Println("начало")
}
// начало
// один
// два
// три
```

Отложенные вызовы выполняются в порядке **LIFO (Last In, First Out)** — стек. Последний `defer` выполняется первым.

Аргументы вычисляются **в момент defer**, а не в момент выполнения:

```go
func main() {
    x := 10
    defer fmt.Println("x =", x)  // x вычисляется СЕЙЧАС = 10
    x = 20
    fmt.Println("main exit")
}
// main exit
// x = 10  — не 20!
```

Это важное правило: значения аргументов зафиксированы при `defer`, но если передаётся **указатель** или используется **замыкание** — увидишь актуальное значение:

```go
func main() {
    x := 10
    defer func() {
        fmt.Println("x =", x)  // замыкание — видит актуальный x
    }()
    x = 20
}
// x = 20
```

---

## Основные паттерны использования

### Закрытие файлов

```go
func readFile(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()  // гарантированно закроется при выходе

    return io.ReadAll(f)
}
```

`defer f.Close()` — сразу после открытия. Так ты никогда не забудешь закрыть файл, даже если функция вернётся раньше из-за ошибки.

### Разблокировка mutex

```go
var mu sync.Mutex
var balance float64

func withdraw(amount float64) error {
    mu.Lock()
    defer mu.Unlock()  // разблокируется при любом выходе

    if balance < amount {
        return errors.New("недостаточно средств")
    }
    balance -= amount
    return nil
}
```

Без `defer` нужно помнить вызывать `mu.Unlock()` в каждой точке выхода из функции — это легко забыть.

### Закрытие соединения с БД

```go
func queryUser(db *sql.DB, id int) (*User, error) {
    rows, err := db.Query("SELECT * FROM users WHERE id = ?", id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // закрыть курсор

    if rows.Next() {
        var u User
        if err := rows.Scan(&u.Name, &u.Email); err != nil {
            return nil, err
        }
        return &u, nil
    }
    return nil, sql.ErrNoRows
}
```

### Измерение времени выполнения

```go
func timeTrack(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s выполнялась %v\n", name, time.Since(start))
    }
}

func expensiveOperation() {
    defer timeTrack("expensiveOperation")()
    // ... долгие вычисления ...
    time.Sleep(100 * time.Millisecond)
}

// expensiveOperation выполнялась 100.4ms
```

Обрати внимание на `defer f()()` — первые скобки вызывают `timeTrack` (сразу, фиксируя время старта), вторые скобки вызывают возвращённую функцию (при defer).

---

## defer + panic + recover

Когда в Go происходит `panic`, стек раскручивается и все отложенные функции выполняются. `recover()` внутри `defer` может перехватить панику:

```go
func safeOperation() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("паника: %v", r)
        }
    }()

    // Потенциально паникующий код:
    var s []int
    _ = s[10]  // panic: index out of range

    return nil
}

func main() {
    err := safeOperation()
    fmt.Println(err)  // паника: runtime error: index out of range [10] with length 0
}
```

**Важные правила recover:**
1. `recover()` работает только **внутри defer**
2. `recover()` останавливает панику только в **текущей горутине**
3. Возвращает значение, переданное в `panic()`, или `nil` если паники не было

```go
func main() {
    // recover без defer — бесполезен:
    r := recover()
    fmt.Println(r)  // всегда nil

    // Правильно — только в defer:
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("восстановились от паники:", r)
        }
    }()

    panic("что-то пошло не так")
}
```

---

## Изменение именованных возвращаемых значений через defer

`defer` может изменить именованные возвращаемые значения:

```go
func mustPositive(x int) (result int, err error) {
    defer func() {
        if result < 0 {
            err = fmt.Errorf("отрицательный результат: %d", result)
            result = 0
        }
    }()

    result = x * x - 10  // может быть отрицательным
    return
}

fmt.Println(mustPositive(2))  // 0, отрицательный результат: -6
fmt.Println(mustPositive(5))  // 15, nil
```

---

## Ловушка: defer в цикле

Из урока о циклах повторим: `defer` выполняется при выходе из **функции**, не из итерации цикла.

```go
// ПЛОХО: все файлы открыты до выхода из функции
func badProcessFiles(paths []string) {
    for _, path := range paths {
        f, _ := os.Open(path)
        defer f.Close()  // накапливаются, закроются только в конце функции!
        process(f)
    }
}

// ХОРОШО: вынести в отдельную функцию
func processOneFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // закроется при выходе из processOneFile
    return process(f)
}

func goodProcessFiles(paths []string) {
    for _, path := range paths {
        processOneFile(path)
    }
}
```

---

## Практический пример: HTTP-обработчик с cleanup

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Логирование начала и конца запроса:
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("%s %s — %v", r.Method, r.URL.Path, duration)
    }()

    // Работа с транзакцией БД:
    tx, err := db.Begin()
    if err != nil {
        http.Error(w, "ошибка БД", 500)
        return
    }
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
        tx.Commit()
    }()

    // Логика обработки...
    if err = processRequest(tx, r); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    w.WriteHeader(200)
}
```

---

## Итог

- `defer` откладывает вызов до выхода из функции
- Порядок выполнения LIFO — последний defer выполняется первым
- Аргументы вычисляются в момент defer, но замыкания видят актуальные значения
- Основные применения: закрытие файлов, разблокировка mutex, логирование, rollback транзакций
- `defer + recover()` — способ перехватить панику
- В циклах — выноси тело в отдельную функцию, иначе ресурсы накапливаются
