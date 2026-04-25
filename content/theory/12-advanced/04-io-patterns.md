---
title: "Паттерны I/O"
description: "io.Reader/Writer композиция, bufio, TeeReader, MultiWriter, потоковая обработка"
order: 4
---

# Паттерны I/O

Go строит всё вокруг двух интерфейсов: `io.Reader` и `io.Writer`. Мы уже видели их в главе об интерфейсах. Здесь — продвинутые паттерны компоновки.

## Интерфейсы-основы

```go
type Reader interface {
    Read(p []byte) error (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Комбинации из пакета io:
type ReadWriter interface { Reader; Writer }
type ReadCloser interface { Reader; Closer }
type WriteCloser interface { Writer; Closer }
type ReadWriteCloser interface { Reader; Writer; Closer }
```

Всё что реализует эти интерфейсы — файлы, HTTP-тело, bytes.Buffer, strings.Reader — взаимозаменяемо.

---

## strings.NewReader и bytes.Buffer

```go
// strings.NewReader — читать из строки
r := strings.NewReader("hello world")
data, _ := io.ReadAll(r)
fmt.Println(string(data))  // hello world

// bytes.Buffer — чтение И запись из памяти
var buf bytes.Buffer
buf.WriteString("hello")
buf.WriteByte(' ')
fmt.Fprintf(&buf, "world %d", 42)
fmt.Println(buf.String())  // hello world 42

// bytes.Buffer как io.Reader:
io.Copy(os.Stdout, &buf)
```

---

## io.TeeReader — чтение с дублированием

`TeeReader` при чтении копирует байты ещё и в Writer:

```go
// Логировать HTTP-тело не потребляя его:
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var bodyBuf bytes.Buffer
        // tee: читать из r.Body и одновременно писать в bodyBuf
        tee := io.TeeReader(r.Body, &bodyBuf)
        r.Body = io.NopCloser(tee)

        next.ServeHTTP(w, r)

        // теперь bodyBuf содержит тело запроса для логирования
        log.Printf("request body: %s", bodyBuf.String())
    })
}
```

```go
// Читать файл и одновременно считать контрольную сумму:
import (
    "crypto/sha256"
    "io"
    "os"
)

func hashFile(path string) ([]byte, []byte, error) {
    f, err := os.Open(path)
    if err != nil { return nil, nil, err }
    defer f.Close()

    h := sha256.New()
    tee := io.TeeReader(f, h)  // f → tee → h (хеш считается попутно)

    data, err := io.ReadAll(tee)
    if err != nil { return nil, nil, err }

    return data, h.Sum(nil), nil
}
```

---

## io.MultiWriter — запись в несколько мест

```go
// Записать одновременно в файл и stdout:
file, _ := os.Create("output.log")
defer file.Close()

multi := io.MultiWriter(os.Stdout, file)
fmt.Fprintln(multi, "это появится в консоли и в файле")
```

Практическое применение — логирование с дублированием:

```go
func setupLogger(logFile string) (*log.Logger, error) {
    f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil { return nil, err }

    multi := io.MultiWriter(os.Stderr, f)
    return log.New(multi, "", log.LstdFlags), nil
}
```

---

## bufio — буферизованный I/O

Буферизация уменьшает количество системных вызовов:

```go
// Запись без буфера: каждый Fprintln = системный вызов write
file, _ := os.Create("out.txt")
for i := 0; i < 10000; i++ {
    fmt.Fprintln(file, i)  // 10000 системных вызовов
}

// Запись с буфером: данные накапливаются, flush редко
file, _ = os.Create("out_buffered.txt")
w := bufio.NewWriterSize(file, 64*1024)  // 64KB буфер
for i := 0; i < 10000; i++ {
    fmt.Fprintln(w, i)
}
w.Flush()  // сбросить остаток буфера — обязательно!
```

Буферизованное чтение построчно:

```go
func processLines(r io.Reader) error {
    scanner := bufio.NewScanner(r)

    // Увеличить буфер для длинных строк:
    buf := make([]byte, 1024*1024)
    scanner.Buffer(buf, len(buf))

    for scanner.Scan() {
        line := scanner.Text()
        // обрабатываем строку
        _ = line
    }
    return scanner.Err()
}
```

---

## Потоковая обработка больших файлов

Никогда не загружай в память файл целиком если он может быть большим:

```go
// НЕВЕРНО: загружаем весь файл
data, err := os.ReadFile("huge.csv")  // может быть несколько GB
records := parse(data)

// ВЕРНО: потоковая обработка
file, err := os.Open("huge.csv")
if err != nil { return err }
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    processLine(scanner.Text())  // обрабатываем строку за строкой
}
return scanner.Err()
```

### Обработка большого JSON-файла

```go
// json.Decoder для потокового чтения массива:
func processJSONStream(r io.Reader) error {
    dec := json.NewDecoder(r)

    // Читаем открывающую скобку массива [
    if _, err := dec.Token(); err != nil {
        return err
    }

    for dec.More() {
        var item Item
        if err := dec.Decode(&item); err != nil {
            return err
        }
        processItem(item)
    }

    return nil
}
```

---

## io.Pipe — соединить reader и writer

`io.Pipe` создаёт синхронный канал между двумя горутинами:

```go
func main() {
    pr, pw := io.Pipe()

    // Горутина-писатель: архивирует данные
    go func() {
        defer pw.Close()
        gw := gzip.NewWriter(pw)
        defer gw.Close()
        fmt.Fprintln(gw, "содержимое файла")
    }()

    // Горутина-читатель: читает сжатые данные
    gr, _ := gzip.NewReader(pr)
    defer gr.Close()
    io.Copy(os.Stdout, gr)
}
```

---

## LimitReader — ограничить чтение

Защита от злоумышленников, загружающих гигабайты:

```go
func handleUpload(w http.ResponseWriter, r *http.Request) {
    // Ограничить тело запроса 10MB
    r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

    data, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "тело слишком большое", http.StatusRequestEntityTooLarge)
        return
    }
    // обработка data
}

// Или через io.LimitReader:
limited := io.LimitReader(r.Body, 10<<20)
data, err := io.ReadAll(limited)
```

---

## Итог

| Инструмент | Назначение |
|-----------|-----------|
| `io.TeeReader(r, w)` | Читать из r, одновременно писать в w |
| `io.MultiWriter(w1, w2, ...)` | Писать во все writers одновременно |
| `bufio.NewWriter(w)` | Буферизованная запись; не забывай Flush() |
| `bufio.NewScanner(r)` | Построчное чтение |
| `io.Pipe()` | Синхронный канал между горутинами |
| `io.LimitReader(r, n)` | Читать не более n байт |
| `json.NewDecoder(r)` | Потоковый JSON без загрузки в память |

Центральный принцип: функции принимают `io.Reader`/`io.Writer`, не конкретные типы — это делает их универсальными и тестируемыми.
