---
title: "Комментарии и именование"
description: "Соглашения Go: camelCase, exported names, godoc-комментарии и go vet"
order: 5
---

# Комментарии и именование

Go — язык с сильными соглашениями. Именование и комментарии в Go — не вопрос вкуса, а часть языковой экосистемы. Правильно именованный код автоматически документируется через `go doc`, а нарушения соглашений приводят к реальным проблемам.

## Соглашения об именовании

### camelCase для всего

Go использует **camelCase**, а не snake_case:

```go
// Правильно:
userName := "alice"
maxRetries := 3
isActive := true
getUserByID(id int)

// Неправильно (не идиоматично для Go):
user_name := "alice"
max_retries := 3
is_active := true
get_user_by_id(id int)
```

Аббревиатуры пишутся заглавными буквами целиком:

```go
// Правильно (неэкспортированные):
userID   // не userId
parseURL // не parseUrl
htmlBody // не htmlBody — аббревиатура целиком строчная, т.к. идентификатор неэкспортированный
dbURL    // не dbUrl

// Правильно (экспортированные):
UserID
ParseURL
HTMLBody // не HtmlBody
DBURL    // не DbUrl

// В экспортированных именах — особенно важно:
type UserID int     // не UserId
func ParseURL(s string) *URL  // не ParseUrl
```

### Экспортированные имена (Exported)

Если имя начинается с **заглавной буквы** — оно экспортировано (видно из других пакетов). Строчная буква — только внутри пакета:

```go
package user

// Экспортированы — доступны из других пакетов:
type User struct {
    Name  string  // экспортированное поле
    Email string
    age   int     // приватное поле
}

func NewUser(name, email string) *User {...}  // экспортированная функция
func (u *User) Validate() error {...}         // экспортированный метод

// Не экспортированы — только внутри пакета user:
func validateEmail(email string) bool {...}
var defaultTimeout = 30 * time.Second
```

**Это не просто конвенция — это часть языка.** Компилятор Go обеспечивает инкапсуляцию через первую букву имени.

### Длина имён

Go придерживается **коротких имён** для малого scope и более длинных для широкого:

```go
// Короткие переменные уместны в маленьком scope:
for i := 0; i < len(items); i++ {...}
for _, v := range values {...}

// В сигнатурах функций — чуть длиннее:
func sum(nums []int) int {...}

// На уровне пакета — длиннее и описательнее:
var userSessionTimeout = 30 * time.Minute
```

Плохой знак — аббревиатуры там, где лучше полное слово:

```go
// Плохо:
func getPrdsLstFrmDB() []Product {...}

// Хорошо:
func getProductsFromDatabase() []Product {...}
```

### Именование интерфейсов

Однометодные интерфейсы называются по имени метода + суффикс `-er`:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Stringer interface {
    String() string
}

type Closer interface {
    Close() error
}
```

### Именование пакетов

- Строчные буквы, одно слово без подчёркиваний
- Краткое, но понятное
- Не используй `util`, `common`, `misc` — это антипаттерн

```go
// Хорошо:
package user
package auth
package store
package http  // уже занято стандартной библиотекой

// Плохо:
package userUtils
package common
package myPackage
```

---

## Комментарии

Go поддерживает два вида комментариев:

```go
// Однострочный комментарий

/*
   Многострочный комментарий
   (используется редко)
*/
```

### godoc-комментарии

Go автоматически генерирует документацию из комментариев. Инструмент `go doc` (и `godoc`) читает комментарии прямо перед объявлением.

**Правила godoc-комментариев:**

1. Комментарий должен быть **прямо перед объявлением** (без пустой строки)
2. Начинается с имени того, что документируется
3. Полное предложение с заглавной буквы и точкой в конце

```go
// User представляет пользователя системы.
type User struct {
    Name  string
    Email string
}

// NewUser создаёт нового пользователя с заданными именем и email.
// Возвращает ошибку, если email имеет неверный формат.
func NewUser(name, email string) (*User, error) {
    if !isValidEmail(email) {
        return nil, fmt.Errorf("неверный формат email: %s", email)
    }
    return &User{Name: name, Email: email}, nil
}

// Validate проверяет корректность данных пользователя.
func (u *User) Validate() error {
    if u.Name == "" {
        return errors.New("имя пользователя не может быть пустым")
    }
    return nil
}
```

Просмотр документации:

```bash
go doc .               # документация текущего пакета
go doc User            # документация типа User
go doc User.Validate   # документация метода
go doc fmt.Printf      # документация из стандартной библиотеки
```

### Комментарии к пакету

Пакет можно документировать комментарием перед строкой `package` — особенно это полезно для публичных пакетов:

```go
// Package user предоставляет функциональность для управления
// пользователями: создание, аутентификация, управление профилем.
//
// Пример использования:
//
//	u, err := user.NewUser("alice", "alice@example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
package user
```

Для больших пакетов комментарий выносят в отдельный файл `doc.go`.

### Когда писать комментарии

Go-сообщество придерживается правила: **комментарий должен объяснять "почему", а не "что"**. Хорошо именованный код сам объясняет что он делает.

```go
// Плохой комментарий — объясняет очевидное:
// Увеличить счётчик на 1
count++

// Хороший комментарий — объясняет неочевидное:
// Используем Atoi вместо ParseInt для совместимости с legacy API,
// который возвращает числа без знака типа, но как string.
n, _ := strconv.Atoi(legacyValue)
```

---

## go vet — анализатор кода

`go vet` — встроенный статический анализатор, который ловит распространённые ошибки:

```bash
go vet ./...
```

Что проверяет:

```go
// Неправильное число аргументов в Printf:
fmt.Printf("%d %d", 1)  // vet: missing argument for fmt verb %d

// Unreachable code:
return
fmt.Println("никогда не выполнится")  // vet: unreachable code

// Неверное использование sync.Mutex:
var mu sync.Mutex
mu2 := mu  // vet: assignment copies lock value

// Shadowing err в тестах:
// (и многие другие проверки)
```

IDE подсвечивает большинство из этих ошибок прямо в редакторе, но `go vet ./...` удобен в CI или перед коммитом.

### gofmt — форматирование кода

Go имеет официальный форматировщик `gofmt`. Нет споров о стиле — весь код форматируется одинаково:

```bash
gofmt -w .       # форматировать все .go файлы
gofmt -d main.go # показать diff без изменения файла
```

В большинстве IDE это происходит автоматически при сохранении.

`gofmt` — стандарт de facto в Go-сообществе, а не рекомендация.

---

## Полный пример: хорошо оформленный пакет

```go
// Package calculator предоставляет базовые математические операции.
package calculator

import "errors"

// ErrDivisionByZero возвращается при попытке деления на ноль.
var ErrDivisionByZero = errors.New("деление на ноль")

// Add возвращает сумму двух чисел.
func Add(a, b float64) float64 {
    return a + b
}

// Divide делит a на b и возвращает результат.
// Возвращает ErrDivisionByZero, если b равно нулю.
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, ErrDivisionByZero
    }
    return a / b, nil
}

// Operation представляет математическую операцию над двумя числами.
type Operation func(a, b float64) float64

// Apply применяет операцию op к числам a и b.
func Apply(op Operation, a, b float64) float64 {
    return op(a, b)
}
```

---

## Итог

**Именование:**
- camelCase (не snake_case)
- Заглавная буква = экспортировано, строчная = приватно
- Аббревиатуры полностью заглавными: `userID`, `parseURL`
- Короткие имена для маленького scope, длинные — для широкого

**Комментарии:**
- Начинаются с имени документируемого объекта
- Полные предложения с точкой
- Объясняют "почему", а не "что"

**Инструменты:**
- `gofmt` — обязательное форматирование
- `go vet` — поиск распространённых ошибок
- `go doc` — просмотр документации
