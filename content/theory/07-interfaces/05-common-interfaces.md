---
title: "Стандартные интерфейсы"
description: "io.Reader/Writer/Closer, fmt.Stringer, error, sort.Interface — реализация своих типов"
order: 5
---

# Стандартные интерфейсы

Go имеет небольшой набор интерфейсов в стандартной библиотеке, которые используются повсеместно. Реализация этих интерфейсов делает твои типы совместимыми со всей экосистемой Go.

## io.Reader

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

Читает данные в буфер `p`, возвращает количество прочитанных байт и ошибку. При достижении конца данных возвращает `io.EOF`.

```go
// Реализация собственного Reader:
type StringReader struct {
    data string
    pos  int
}

func (r *StringReader) Read(p []byte) (n int, err error) {
    if r.pos >= len(r.data) {
        return 0, io.EOF
    }
    n = copy(p, r.data[r.pos:])
    r.pos += n
    return n, nil
}

// Работает с любой функцией, принимающей io.Reader:
reader := &StringReader{data: "Hello, World!"}
buf := make([]byte, 4)
for {
    n, err := reader.Read(buf)
    if n > 0 {
        fmt.Print(string(buf[:n]))
    }
    if err == io.EOF {
        break
    }
}
// Hello, World!
```

Всё, что принимает `io.Reader`, работает с файлами, HTTP-ответами, строками (`strings.NewReader`), байтовыми буферами и любым твоим типом.

---

## io.Writer

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

Записывает байты из `p`, возвращает количество записанных и ошибку.

```go
// Подсчёт записанных байт:
type CountingWriter struct {
    w     io.Writer
    count int64
}

func (cw *CountingWriter) Write(p []byte) (n int, err error) {
    n, err = cw.w.Write(p)
    cw.count += int64(n)
    return
}

// Использование:
var buf bytes.Buffer
cw := &CountingWriter{w: &buf}

fmt.Fprintf(cw, "Hello, %s!\n", "World")
fmt.Fprintf(cw, "Байт записано: %d\n", cw.count)

fmt.Println(buf.String())
// Hello, World!
// Байт записано: 14
```

---

## io.Closer

```go
type Closer interface {
    Close() error
}
```

Освобождение ресурсов. Часто встречается в составных интерфейсах:

```go
type ReadCloser interface {
    Reader
    Closer
}

type WriteCloser interface {
    Writer
    Closer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

Пример: `http.Response.Body` имеет тип `io.ReadCloser` — нужно обязательно закрывать:

```go
resp, err := http.Get("https://example.com")
if err != nil {
    return err
}
defer resp.Body.Close()  // обязательно!

body, err := io.ReadAll(resp.Body)
```

---

## fmt.Stringer

```go
type Stringer interface {
    String() string
}
```

Реализуй `String()`, чтобы управлять тем, как тип выводится через `fmt.Println`, `fmt.Printf("%v")` и т.д.:

```go
type Duration struct {
    hours   int
    minutes int
    seconds int
}

func (d Duration) String() string {
    return fmt.Sprintf("%02d:%02d:%02d", d.hours, d.minutes, d.seconds)
}

d := Duration{1, 30, 45}
fmt.Println(d)        // 01:30:45
fmt.Printf("%v\n", d) // 01:30:45
s := fmt.Sprint(d)    // "01:30:45"
```

---

## error

```go
type error interface {
    Error() string
}
```

Уже знакомый интерфейс. Реализуй для создания кастомных ошибок:

```go
type HTTPError struct {
    Code    int
    Message string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

func fetch(url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return &HTTPError{
            Code:    resp.StatusCode,
            Message: resp.Status,
        }
    }
    return nil
}

err := fetch("https://example.com/api")
if he, ok := err.(*HTTPError); ok {
    fmt.Printf("HTTP ошибка %d: %s\n", he.Code, he.Message)
}
```

---

## sort.Interface

```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}
```

Реализуй для кастомной сортировки любого типа:

```go
type Person struct {
    Name string
    Age  int
}

// Сортировка по имени:
type ByName []Person

func (s ByName) Len() int           { return len(s) }
func (s ByName) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s ByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Сортировка по возрасту:
type ByAge []Person

func (s ByAge) Len() int           { return len(s) }
func (s ByAge) Less(i, j int) bool { return s[i].Age < s[j].Age }
func (s ByAge) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func main() {
    people := []Person{
        {"Карл", 35},
        {"Алиса", 25},
        {"Боб", 30},
    }

    sort.Sort(ByAge(people))
    fmt.Println(people)  // [{Алиса 25} {Боб 30} {Карл 35}]

    sort.Sort(ByName(people))
    fmt.Println(people)  // [{Алиса 25} {Боб 30} {Карл 35}]
}
```

На практике проще использовать `sort.Slice` — реализовывать `sort.Interface` нужно только если хочешь переиспользуемый тип сортировки.

---

## Реализация нескольких интерфейсов

Один тип может реализовывать сколько угодно интерфейсов:

```go
type Buffer struct {
    data []byte
    pos  int
}

// Реализует io.Reader:
func (b *Buffer) Read(p []byte) (n int, err error) {
    if b.pos >= len(b.data) {
        return 0, io.EOF
    }
    n = copy(p, b.data[b.pos:])
    b.pos += n
    return n, nil
}

// Реализует io.Writer:
func (b *Buffer) Write(p []byte) (n int, err error) {
    b.data = append(b.data, p...)
    return len(p), nil
}

// Реализует fmt.Stringer:
func (b *Buffer) String() string {
    return string(b.data)
}

// Реализует io.ReadWriter, io.Reader, io.Writer, fmt.Stringer — всё сразу
var buf Buffer
fmt.Fprintf(&buf, "Hello, %s!", "World")  // использует io.Writer
fmt.Println(&buf)                          // использует fmt.Stringer
```

---

## Итог

| Интерфейс | Применение |
|-----------|-----------|
| `io.Reader` | Любой источник данных: файлы, HTTP, строки, сетевые соединения |
| `io.Writer` | Любой приёмник: файлы, буферы, stdout, сетевые соединения |
| `io.Closer` | Освобождение ресурсов; всегда вызывай через `defer` |
| `fmt.Stringer` | Контроль текстового представления типа |
| `error` | Кастомные ошибки с дополнительными данными |
| `sort.Interface` | Кастомная сортировка коллекций |

Реализация этих интерфейсов делает твои типы совместимыми со стандартной библиотекой и тысячами сторонних пакетов.
