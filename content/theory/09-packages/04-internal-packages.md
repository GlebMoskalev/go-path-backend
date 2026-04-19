---
title: "Internal-пакеты и структура проекта"
description: "internal/, структура большого проекта, циклические зависимости"
order: 4
---

# Internal-пакеты и структура проекта

По мере роста проекта важно правильно организовать код. Go предоставляет механизм `internal/` для строгого контроля видимости пакетов.

## Пакет internal/

Директория `internal/` — специальная директива компилятора. Пакеты внутри неё доступны **только из родительской директории и её поддиректорий**:

```
myapp/
├── main.go              # может импортировать internal/
├── cmd/
│   └── server/
│       └── main.go      # может импортировать internal/
├── internal/
│   ├── auth/            # ТОЛЬКО для myapp, никто снаружи
│   │   └── auth.go
│   └── db/
│       └── db.go
└── pkg/
    └── models/          # публичный пакет — для всех
        └── models.go
```

```go
// Файл: myapp/main.go — OK
import "myapp/internal/auth"

// Файл: другой-модуль/main.go — ОШИБКА КОМПИЛЯЦИИ
import "myapp/internal/auth"
// use of internal package myapp/internal/auth not allowed
```

Правило доступа: пакет `a/b/c/internal/d/e/f` доступен из `a/b/c` и всего, что внутри него.

---

## Структура реального Go-проекта

Нет единого «правильного» шаблона, но вот наиболее распространённый:

```
myapp/
├── cmd/                    # исполняемые команды
│   ├── server/
│   │   └── main.go
│   └── migrate/
│       └── main.go
├── internal/               # приватные пакеты
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── storage/
│   │   ├── postgres.go
│   │   └── redis.go
│   └── config/
│       └── config.go
├── pkg/                    # публичные пакеты (для внешнего использования)
│   ├── models/
│   │   └── user.go
│   └── api/
│       └── client.go
├── api/                    # OpenAPI/Protobuf определения
│   └── v1/
│       └── service.proto
├── migrations/             # SQL-миграции
├── scripts/                # вспомогательные скрипты
├── go.mod
├── go.sum
└── README.md
```

**Минимальная структура** для небольшого проекта:

```
myapp/
├── main.go
├── handler.go
├── service.go
├── storage.go
├── go.mod
└── go.sum
```

Не усложняй структуру заранее. Начни с плоской структуры, вводи пакеты по мере необходимости.

---

## Циклические зависимости

Go **запрещает** циклические зависимости между пакетами:

```
Пакет A импортирует B
Пакет B импортирует A
→ ОШИБКА: import cycle not allowed
```

```go
// package a
import "myapp/b"  // A зависит от B

// package b
import "myapp/a"  // B зависит от A — ЦИКЛ!
```

Компилятор Go выдаст ошибку: `import cycle not allowed`.

### Как избежать циклов

**Стратегия 1: Выделить общие типы в отдельный пакет**

```
Было:           Стало:
a ↔ b           a → types ← b
                a → types
                b → types
```

```go
// package types — только типы, никаких зависимостей
package types

type User struct {
    ID   int
    Name string
}

// package user — логика пользователей
import "myapp/types"

// package order — логика заказов
import "myapp/types"
// order и user теперь не зависят друг от друга!
```

**Стратегия 2: Использовать интерфейсы**

```go
// Пакет service хочет использовать repository, но не создавать цикл:
package service

// Определяем интерфейс здесь, а не импортируем конкретный тип:
type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(u *User) error
}

type UserService struct {
    repo UserRepository  // интерфейс, не конкретный тип
}
```

**Стратегия 3: Слои (dependency direction)**

Следи за направлением зависимостей. Типичный порядок:

```
HTTP Handlers → Services → Repository → Database
     ↓               ↓           ↓
   Models          Models      Models
```

Каждый слой зависит только от нижeleжащих. Никаких зависимостей наверх.

---

## pkg/ vs internal/

| | `internal/` | `pkg/` |
|---|---|---|
| Видимость | Только из родительского модуля | Всем, включая внешние модули |
| Назначение | Приватная реализация | Публичный API для внешних пользователей |
| Когда использовать | Логика, которую не хочешь экспортировать | Библиотека для использования другими |

Если ты пишешь приложение (не библиотеку) — весь код можно класть в `internal/`. Если пишешь библиотеку — используй `pkg/` для публичного API.

---

## Практический пример: разделение пакетов

```
users/
├── cmd/
│   └── main.go
├── internal/
│   ├── handler/       # HTTP обработчики
│   │   └── user.go
│   ├── service/       # бизнес-логика
│   │   └── user.go
│   ├── repository/    # работа с данными
│   │   └── user.go
│   └── domain/        # доменные типы (общие для всех)
│       └── user.go
└── go.mod
```

```go
// domain/user.go — только типы, нет зависимостей
package domain

type User struct {
    ID    int
    Name  string
    Email string
}

type UserRepository interface {
    FindByID(id int) (*User, error)
}
```

```go
// repository/user.go — реализация репозитория
package repository

import (
    "myapp/internal/domain"
)

type PostgresUserRepo struct{ db *sql.DB }

func (r *PostgresUserRepo) FindByID(id int) (*domain.User, error) {
    // ...
}
```

```go
// service/user.go — бизнес-логика
package service

import (
    "myapp/internal/domain"
)

type UserService struct {
    repo domain.UserRepository  // зависим от интерфейса, не от реализации
}
```

---

## Итог

- `internal/` — пакеты доступны только из родительской директории
- Используй `internal/` для приватной реализации
- Циклические зависимости запрещены — решай через общие типы или интерфейсы
- Структура должна отражать направление зависимостей (снаружи → внутрь)
- Не создавай сложную структуру заранее — начинай просто
