---
title: "Стандартная библиотека"
description: "Обзор ключевых пакетов: os, io, bufio, path/filepath, time, math, sort, strings, strconv, regexp"
order: 3
---

# Стандартная библиотека

Go поставляется с богатой стандартной библиотекой. Перед тем как искать внешние пакеты — проверь, нет ли готового решения в stdlib.

## os — работа с операционной системой

```go
import "os"

// Файлы:
f, err := os.Open("input.txt")    // только чтение
f, err := os.Create("output.txt") // создать/перезаписать
f, err := os.OpenFile("file.txt", os.O_APPEND|os.O_WRONLY, 0644)

defer f.Close()

// Директории:
os.Mkdir("mydir", 0755)
os.MkdirAll("a/b/c", 0755)     // создать все промежуточные
os.Remove("file.txt")
os.RemoveAll("dir/")
os.Rename("old.txt", "new.txt")

// Информация о файле:
info, err := os.Stat("file.txt")
if os.IsNotExist(err) {
    fmt.Println("файл не существует")
}
if err == nil {
    fmt.Println("Размер:", info.Size())
    fmt.Println("Изменён:", info.ModTime())
}

// Переменные окружения:
home := os.Getenv("HOME")
os.Setenv("MY_VAR", "value")

// Аргументы командной строки:
fmt.Println(os.Args)  // ["/path/to/binary", "arg1", "arg2"]

// Стандартные потоки:
os.Stdout, os.Stderr, os.Stdin
```

## io — абстрактный ввод-вывод

```go
import "io"

// Чтение всего содержимого:
data, err := io.ReadAll(reader)

// Копирование:
n, err := io.Copy(dst, src)

// Ограниченное чтение:
limited := io.LimitReader(reader, 1024)  // не более 1024 байт

// Мультиплексирование записи:
multi := io.MultiWriter(file1, file2, os.Stdout)

// Отбрасывание (devnull):
io.Discard
```

## bufio — буферизованный ввод-вывод

```go
import "bufio"

// Построчное чтение файла (самый частый паттерн):
file, _ := os.Open("file.txt")
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}
if err := scanner.Err(); err != nil {
    log.Fatal(err)
}

// Буферизованная запись:
writer := bufio.NewWriter(file)
fmt.Fprintln(writer, "строка 1")
fmt.Fprintln(writer, "строка 2")
writer.Flush()  // не забывай сбросить буфер!

// Буферизованное чтение:
reader := bufio.NewReader(os.Stdin)
line, _ := reader.ReadString('\n')
```

## path/filepath — работа с путями файловой системы

```go
import "path/filepath"

// Построение путей (OS-зависимо: / на Unix, \ на Windows):
p := filepath.Join("usr", "local", "bin")   // usr/local/bin

// Разбор пути:
filepath.Dir("/usr/local/bin/go")   // /usr/local/bin
filepath.Base("/usr/local/bin/go")  // go
filepath.Ext("file.go")             // .go

// Абсолютный путь:
abs, _ := filepath.Abs("relative/path")

// Glob — поиск файлов по шаблону:
files, _ := filepath.Glob("*.go")

// Обход директории:
filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
    if err != nil { return err }
    fmt.Println(path)
    return nil
})

// Более новый вариант с filepath.WalkDir (эффективнее):
filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
    if d.IsDir() { return nil }
    fmt.Println(path)
    return nil
})
```

## time — работа со временем

```go
import "time"

// Текущее время:
now := time.Now()
fmt.Println(now)                             // 2026-04-19 14:32:01.123456 +0300 MSK

// Форматирование — используется магическая дата: Mon Jan 2 15:04:05 MST 2006
now.Format("2006-01-02")           // 2026-04-19
now.Format("15:04:05")             // 14:32:01
now.Format("02.01.2006 15:04:05") // 19.04.2026 14:32:01
now.Format(time.RFC3339)           // 2026-04-19T14:32:01+03:00

// Парсинг:
t, _ := time.Parse("2006-01-02", "2026-04-19")
t, _ := time.Parse(time.RFC3339, "2026-04-19T14:32:01+03:00")

// Длительности:
d := 2*time.Hour + 30*time.Minute
time.Sleep(100 * time.Millisecond)

// Арифметика:
future := now.Add(24 * time.Hour)
duration := future.Sub(now)        // 24h0m0s

// Сравнение:
now.Before(future)  // true
now.After(future)   // false

// Unix timestamp:
now.Unix()      // секунды с 1970
now.UnixMilli() // миллисекунды
now.UnixNano()  // наносекунды

// Таймер и тикер:
timer := time.NewTimer(5 * time.Second)
<-timer.C  // ждать 5 секунд

ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for tick := range ticker.C {
    fmt.Println(tick)  // каждую секунду
}
```

## math — математические функции

```go
import "math"

math.Abs(-5.0)          // 5.0
math.Sqrt(16.0)         // 4.0
math.Pow(2, 10)         // 1024.0
math.Log(math.E)        // 1.0
math.Log2(1024)         // 10.0
math.Log10(100)         // 2.0
math.Ceil(1.2)          // 2.0
math.Floor(1.8)         // 1.0
math.Round(1.5)         // 2.0
math.Min(1.0, 2.0)      // 1.0
math.Max(1.0, 2.0)      // 2.0

math.Pi     // 3.141592653589793
math.E      // 2.718281828459045
math.MaxFloat64
math.MaxInt
```

## sort — сортировка

```go
import "sort"

// Примитивные слайсы:
nums := []int{3, 1, 4, 1, 5}
sort.Ints(nums)
fmt.Println(sort.IntsAreSorted(nums))  // true

strs := []string{"banana", "apple", "cherry"}
sort.Strings(strs)

// Кастомная сортировка:
sort.Slice(nums, func(i, j int) bool { return nums[i] > nums[j] })  // по убыванию
sort.SliceStable(nums, func(i, j int) bool { return nums[i] < nums[j] })

// Бинарный поиск:
i := sort.SearchInts([]int{1, 3, 5, 7, 9}, 5)  // 2
```

## strings — работа со строками

```go
import "strings"

strings.Contains("hello world", "world")   // true
strings.HasPrefix("hello", "hel")          // true
strings.HasSuffix("hello", "llo")          // true
strings.Index("hello", "ll")               // 2
strings.Count("cheese", "e")               // 3

strings.ToUpper("hello")                   // HELLO
strings.ToLower("HELLO")                   // hello
strings.Title("hello world")               // устарело

strings.TrimSpace("  hello  ")             // hello
strings.Trim("***hi***", "*")              // hi
strings.TrimPrefix("go-path", "go-")      // path
strings.TrimSuffix("test.go", ".go")      // test

strings.Split("a,b,c", ",")               // [a b c]
strings.SplitN("a,b,c", ",", 2)           // [a b,c]
strings.Join([]string{"a","b"}, "-")       // a-b
strings.Fields("  foo  bar  ")             // [foo bar]

strings.Replace("oink oink", "oink", "moo", 1)   // moo oink
strings.ReplaceAll("oink oink", "oink", "moo")    // moo moo
strings.Repeat("ab", 3)                    // ababab

// strings.Builder для эффективной конкатенации:
var sb strings.Builder
for i := 0; i < 10; i++ {
    fmt.Fprintf(&sb, "item%d ", i)
}
fmt.Println(sb.String())
```

## regexp — регулярные выражения

```go
import "regexp"

// Компиляция (дорогая операция, делать один раз):
re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

re.MatchString("today is 2026-04-19")  // true
re.FindString("2026-04-19 and 2026-04-20")  // 2026-04-19
re.FindAllString("2026-04-19 and 2026-04-20", -1)  // [2026-04-19 2026-04-20]

// С группами захвата:
re2 := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
matches := re2.FindStringSubmatch("2026-04-19")
// matches[0] = "2026-04-19", [1]="2026", [2]="04", [3]="19"

// Замена:
result := re.ReplaceAllString("date: 2026-04-19", "[DATE]")
// date: [DATE]
```

---

## Итог

Стандартная библиотека покрывает большинство задач без внешних зависимостей:

| Задача | Пакет |
|--------|-------|
| Файловая система | `os`, `io`, `path/filepath` |
| Буферизованный I/O | `bufio` |
| Время и таймеры | `time` |
| Математика | `math`, `math/rand` |
| Строки | `strings`, `strconv`, `unicode/utf8` |
| Сортировка | `sort` |
| Регулярные выражения | `regexp` |
| Сериализация | `encoding/json`, `encoding/xml`, `encoding/csv` |
| HTTP | `net/http`, `net/url` |
| Шаблоны | `text/template`, `html/template` |
| Синхронизация | `sync`, `sync/atomic` |
| Тестирование | `testing` |
