---
title: "Циклы for"
description: "Три формы for, бесконечный цикл, while-аналог, range по string/slice/map/channel"
order: 3
---

# Циклы for

В Go только один вид цикла — `for`. Но он выполняет роль `while`, `do-while`, `foreach` и бесконечного цикла из других языков. Это намеренное упрощение.

## Три формы for

### 1. Классический трёхкомпонентный for (C-style)

```go
for инициализация; условие; постшаг {
    // тело цикла
}
```

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
// 0 1 2 3 4

// Обратный порядок:
for i := 4; i >= 0; i-- {
    fmt.Println(i)
}
// 4 3 2 1 0

// Шаг больше 1:
for i := 0; i <= 100; i += 10 {
    fmt.Print(i, " ")
}
// 0 10 20 30 40 50 60 70 80 90 100
```

Любой из трёх компонентов можно пропустить:

```go
// Только условие (аналог while):
i := 0
for i < 5 {
    fmt.Println(i)
    i++
}

// Без инициализации и постшага (аналог while):
for running {
    process()
}
```

### 2. Цикл с одним условием (while-аналог)

```go
for условие {
    // тело цикла
}
```

```go
n := 1
for n < 1000 {
    n *= 2
}
fmt.Println(n)  // 1024

// Чтение файла построчно:
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}
```

### 3. Бесконечный цикл

```go
for {
    // бесконечный цикл
    // выход только через break, return или panic
}
```

```go
// Типичный server-loop:
for {
    conn, err := listener.Accept()
    if err != nil {
        log.Println("ошибка:", err)
        continue
    }
    go handleConnection(conn)
}

// Цикл с опросом:
for {
    if data, ok := tryRead(); ok {
        process(data)
        break
    }
    time.Sleep(100 * time.Millisecond)
}
```

---

## for range — итерация по коллекциям

`for range` — самая идиоматичная форма. Работает со слайсами, массивами, строками, картами и каналами.

### Range по слайсу/массиву

```go
nums := []int{10, 20, 30, 40, 50}

// Индекс и значение:
for i, v := range nums {
    fmt.Printf("nums[%d] = %d\n", i, v)
}

// Только индекс:
for i := range nums {
    nums[i] *= 2  // изменяем элементы через индекс
}

// Только значение (индекс игнорируем):
for _, v := range nums {
    fmt.Println(v)
}
```

**Важно**: `v` в range — **копия** элемента. Изменение `v` не изменит элемент в слайсе:

```go
nums := []int{1, 2, 3}
for _, v := range nums {
    v *= 10  // меняем копию, nums не изменяется!
}
fmt.Println(nums)  // [1 2 3]

// Правильный способ изменить элементы:
for i := range nums {
    nums[i] *= 10
}
fmt.Println(nums)  // [10 20 30]
```

### Range по строке

```go
s := "Hello, 世界"

for i, r := range s {
    fmt.Printf("позиция %d: '%c' (U+%04X)\n", i, r, r)
}
// позиция 0: 'H' (U+0048)
// позиция 1: 'e' (U+0065)
// ...
// позиция 7: '世' (U+4E16)
// позиция 10: '界' (U+754C)
// Позиции 7 и 10 — байтовые позиции многобайтовых символов
```

Range по строке автоматически декодирует UTF-8 и возвращает руны. Это безопасно для Unicode.

### Range по map

```go
scores := map[string]int{
    "Алиса": 95,
    "Боб":   87,
    "Карл":  92,
}

for name, score := range scores {
    fmt.Printf("%s: %d\n", name, score)
}
// Порядок НЕПРЕДСКАЗУЕМ — карты не упорядочены!
```

```go
// Только ключи:
for name := range scores {
    fmt.Println(name)
}
```

> ⚠️ **Порядок итерации по карте не гарантирован** и намеренно рандомизируется между запусками. Если нужен определённый порядок — собирай ключи в слайс и сортируй.

### Range по каналу

```go
ch := make(chan int)

go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)  // важно: закрыть канал, иначе range зависнет
}()

for v := range ch {
    fmt.Println(v)
}
// 0 1 2 3 4
```

Range по каналу читает значения до тех пор, пока канал не будет закрыт.

> 💡 Каналы подробно разберём в главе 10 «Конкурентность».

---

## break и continue

### break — выход из цикла

```go
for i := 0; i < 100; i++ {
    if i*i > 50 {
        fmt.Println("останавливаемся на", i)
        break
    }
}

// break также работает в switch:
switch x {
case 1:
    if someCondition {
        break  // выходит из switch, не из цикла
    }
    doSomething()
}
```

### continue — следующая итерация

```go
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue  // пропускаем чётные
    }
    fmt.Println(i)  // 1 3 5 7 9
}
```

---

## Типичные ошибки

**Ошибка 1**: Изменение значения в range-копии.

```go
type Point struct{ X, Y int }
points := []Point{{1, 2}, {3, 4}}

for _, p := range points {
    p.X *= 2  // изменяем копию, points не меняется!
}
fmt.Println(points)  // [{1 2} {3 4}]

// Правильно:
for i := range points {
    points[i].X *= 2
}
fmt.Println(points)  // [{2 2} {6 4}]
```

**Ошибка 2**: Бесконечный цикл из-за незакрытого канала.

```go
ch := make(chan int)
go sendData(ch)  // если sendData не закрывает канал...
for v := range ch {  // ...цикл будет ждать вечно
    fmt.Println(v)
}
```

**Ошибка 3**: Рассчитывать на порядок итерации по map.

```go
m := map[string]int{"b": 2, "a": 1, "c": 3}
for k, v := range m {
    fmt.Println(k, v)  // порядок каждый раз разный!
}
```

---

## Итог

- В Go один цикл — `for`, который выполняет роль всех циклов
- `for i := 0; i < n; i++` — C-style, для подсчёта
- `for условие {}` — while-аналог
- `for {}` — бесконечный цикл
- `for i, v := range коллекция` — итерация по слайсу, строке, map, каналу
- `v` в range — копия, изменение не затрагивает оригинал
- Порядок итерации по map непредсказуем
