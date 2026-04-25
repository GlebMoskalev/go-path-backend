---
title: "Основы тестирования"
description: "testing.T, go test, t.Error vs t.Fatal, t.Helper, организация тестов"
order: 1
---

# Основы тестирования

Go поставляется со встроенным инструментом тестирования. Никаких внешних фреймворков не требуется.

## Структура теста

Файл с тестами: имя оканчивается на `_test.go`. Функция теста: начинается с `Test`, принимает `*testing.T`:

```go
// math.go
package math

func Add(a, b int) int {
    return a + b
}

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("деление на ноль")
    }
    return a / b, nil
}
```

```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d, ожидалось 5", result)
    }
}
```

Тесты живут в том же пакете (доступ к неэкспортированным функциям) или в пакете `foo_test` (только публичный API).

---

## Запуск тестов

```bash
go test ./...                    # все пакеты рекурсивно
go test .                        # текущий пакет
go test -v ./...                 # подробный вывод
go test -run TestAdd ./...       # только тесты, соответствующие regexp
go test -run "TestAdd|TestDiv"   # несколько тестов
go test -count=1 ./...           # отключить кэш результатов
go test -timeout 30s ./...       # таймаут (по умолчанию 10m)
```

Пример вывода `-v`:

```
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
=== RUN   TestDivide
--- PASS: TestDivide (0.00s)
PASS
ok      myapp/math    0.001s
```

---

## t.Error vs t.Fatal

| Метод | Поведение |
|-------|-----------|
| `t.Error(args...)` | Записать ошибку, продолжить тест |
| `t.Errorf(format, args...)` | То же с форматированием |
| `t.Fatal(args...)` | Записать ошибку, **остановить** тест |
| `t.Fatalf(format, args...)` | То же с форматированием |
| `t.Log(args...)` | Лог (виден только при `-v` или провале) |
| `t.Skip(args...)` | Пропустить тест |

```go
func TestDivide(t *testing.T) {
    result, err := Divide(10, 2)
    if err != nil {
        t.Fatalf("неожиданная ошибка: %v", err)
        // Fatal останавливает тест здесь
    }
    // Если Fatal не вызван — продолжаем проверки
    if result != 5 {
        t.Errorf("Divide(10, 2) = %f, ожидалось 5", result)
    }
}
```

Используй `Fatal` когда дальнейшее выполнение теста не имеет смысла (nil pointer, неинициализированный объект).

---

## t.Helper — правильные стек-трейсы

Вспомогательная функция должна вызывать `t.Helper()`, чтобы строка ошибки указывала на вызывающий код, а не внутрь helper:

```go
// БЕЗ t.Helper: строка ошибки указывает на строку внутри assertEqual
func assertEqual(t *testing.T, got, want int) {
    if got != want {
        t.Errorf("got %d, want %d", got, want)  // указывает сюда
    }
}

// С t.Helper: строка ошибки указывает на вызывающий тест
func assertEqual(t *testing.T, got, want int) {
    t.Helper()  // добавляем в начало helper-функции
    if got != want {
        t.Errorf("got %d, want %d", got, want)  // теперь строка в тесте
    }
}

func TestAdd(t *testing.T) {
    assertEqual(t, Add(2, 3), 5)   // ← ошибка укажет на эту строку
    assertEqual(t, Add(-1, 1), 0)
}
```

---

## Тестирование ошибок

```go
func TestDivideByZero(t *testing.T) {
    _, err := Divide(10, 0)
    if err == nil {
        t.Fatal("ожидалась ошибка деления на ноль, но err == nil")
    }
}

// Проверка конкретного типа ошибки через errors.As:
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("поле %s: %s", e.Field, e.Message)
}

func TestValidationError(t *testing.T) {
    err := validate("")
    
    var valErr *ValidationError
    if !errors.As(err, &valErr) {
        t.Fatalf("ожидался ValidationError, получили: %T", err)
    }
    if valErr.Field != "name" {
        t.Errorf("ожидалось поле 'name', получили '%s'", valErr.Field)
    }
}
```

---

## Временные ресурсы и очистка: t.Cleanup

`t.Cleanup` регистрирует функцию, вызываемую после завершения теста:

```go
func TestWithTempFile(t *testing.T) {
    f, err := os.CreateTemp("", "test-*.txt")
    if err != nil {
        t.Fatal(err)
    }
    t.Cleanup(func() {
        os.Remove(f.Name())  // гарантированно удалится после теста
    })

    // работаем с файлом...
    fmt.Fprintln(f, "данные")
    f.Close()
    // ...
}
```

`t.Cleanup` надёжнее `defer` когда cleanup нужен в helper-функциях, не в самом тесте.

---

## Организация тестового кода

### TestMain — глобальная настройка

```go
func TestMain(m *testing.M) {
    // setup перед всеми тестами
    db := setupTestDB()
    
    code := m.Run()  // запуск тестов
    
    // teardown после всех тестов
    db.Close()
    
    os.Exit(code)
}
```

### Файловая структура

```
myapp/
├── service.go
├── service_test.go        # тесты в том же пакете (package service)
├── service_external_test.go  # тесты публичного API (package service_test)
└── testdata/
    ├── input.json
    └── expected.json
```

Директория `testdata/` — стандартное место для тестовых фикстур. Go не включает её в бинарник.

---

## Итог

- Файл `*_test.go`, функции `Test*`, параметр `*testing.T`
- `t.Error`/`t.Errorf` — записать ошибку, продолжить тест
- `t.Fatal`/`t.Fatalf` — записать ошибку, остановить тест
- `t.Helper()` — в начале helper-функций для правильных строк ошибок
- `t.Cleanup(func)` — гарантированная очистка ресурсов
- `go test -v -run TestName ./...` — запуск конкретного теста
