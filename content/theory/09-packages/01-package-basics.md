---
title: "Основы пакетов"
description: "Объявление, именование, exported/unexported, init(), порядок инициализации"
order: 1
---

# Основы пакетов

Пакет — единица организации кода в Go. Все файлы в одной директории принадлежат одному пакету.

## Объявление пакета

Каждый файл начинается с `package имя`:

```go
package main       // исполняемая программа
package user       // библиотечный пакет
package http       // пакет стандартной библиотеки
package httptest   // тестовый вспомогательный пакет
```

**Правила именования:**
- Строчные буквы, без подчёркиваний
- Одно слово (краткое и описательное)
- Не `util`, не `common` — конкретное имя
- Имя пакета ≠ имя директории (но обычно совпадают)

---

## Exported и unexported

В Go нет ключевых слов `public`/`private`. Всё определяется регистром первой буквы:

```go
package user

// Exported — доступны из других пакетов:
type User struct {
    Name  string  // экспортированное поле
    Email string  // экспортированное поле
    age   int     // неэкспортированное поле!
}

var DefaultTimeout = 30 * time.Second  // экспортированная переменная

func NewUser(name, email string) *User { ... }  // экспортированная функция
func (u *User) Validate() error { ... }         // экспортированный метод

// Unexported — только внутри пакета:
type session struct { ... }
var internalCache = make(map[string]*User)
func validateEmail(email string) bool { ... }
```

Попытка обратиться к неэкспортированному идентификатору из другого пакета:

```go
// Пакет main:
import "myapp/user"

u := user.NewUser("Alice", "alice@example.com")
fmt.Println(u.Name)   // OK — экспортировано
fmt.Println(u.age)    // ОШИБКА: u.age undefined (cannot refer to unexported field)
```

---

## Импорт пакетов

```go
import "fmt"
import "os"

// Или сгруппированно (рекомендуется):
import (
    "fmt"
    "os"
    "strings"
    
    "github.com/user/repo/pkg"  // внешний пакет
)
```

### Псевдонимы импорта

```go
import (
    "encoding/json"
    
    // Псевдоним для сокращения:
    yaml "gopkg.in/yaml.v3"
    
    // Псевдоним _ — только побочные эффекты (init):
    _ "database/sql/sqlite3"
    
    // Псевдоним . — импорт в текущее пространство имён (избегай!):
    . "math"  // теперь можно писать Sin(x) вместо math.Sin(x)
)

yaml.Marshal(...)  // используем псевдоним
```

Групповой импорт по соглашению организуется блоками: стандартная библиотека, внешние пакеты, внутренние пакеты:

```go
import (
    "fmt"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "myapp/internal/user"
    "myapp/pkg/config"
)
```

---

## init() — функция инициализации

Каждый пакет может иметь одну или несколько функций `init()`. Они выполняются автоматически при загрузке пакета, до `main()`:

```go
package config

import (
    "log"
    "os"
)

var DatabaseURL string

func init() {
    DatabaseURL = os.Getenv("DATABASE_URL")
    if DatabaseURL == "" {
        log.Fatal("DATABASE_URL не установлена")
    }
}
```

Правила `init()`:
- Нет параметров и возвращаемых значений
- Нельзя вызвать явно
- Можно иметь несколько `init()` в одном файле и пакете
- Выполняются по порядку файлов в алфавитном порядке

---

## Порядок инициализации

Go гарантирует строгий порядок инициализации:

1. Импортированные пакеты инициализируются сначала (рекурсивно)
2. Переменные пакетного уровня в порядке объявления
3. Функции `init()` в порядке объявления (файл за файлом)

```
Пакет A зависит от B, B зависит от C:

C инициализируется первым:
  C: переменные уровня пакета
  C: init()

Затем B:
  B: переменные уровня пакета
  B: init()

Затем A:
  A: переменные уровня пакета
  A: init()

Наконец main():
  main: переменные уровня пакета
  main: init()
  main: main()
```

---

## Практический пример: database driver регистрация

Классический пример использования `init()` — регистрация драйверов:

```go
// Пакет github.com/lib/pq (PostgreSQL driver):
package pq

import "database/sql"

func init() {
    sql.Register("postgres", &Driver{})
}

// В вашем коде — импорт только для побочного эффекта:
import (
    "database/sql"
    _ "github.com/lib/pq"  // регистрирует PostgreSQL driver
)

func main() {
    db, err := sql.Open("postgres", "postgresql://...")
    // ...
}
```

---

## Итог

- Пакет = директория; все файлы директории — один пакет
- Большая буква = экспортировано (публично), строчная = приватно
- `init()` выполняется автоматически при загрузке пакета
- Порядок инициализации: сначала зависимости, затем текущий пакет
- Имя пакета: строчное, одно слово, конкретное (не `util`)
- Псевдоним `_` для импорта только с побочными эффектами
