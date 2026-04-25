---
title: "Go Modules"
description: "go mod init/tidy/vendor, go.mod и go.sum, семантическое версионирование, replace"
order: 2
---

# Go Modules

Go Modules — стандартная система управления зависимостями с Go 1.11. Каждый проект — это модуль с уникальным путём и версией.

## Создание модуля

```bash
mkdir myapp
cd myapp
go mod init github.com/myuser/myapp
```

Создаётся файл `go.mod`:

```go
module github.com/myuser/myapp

go 1.26

require (
    // зависимости появятся здесь
)
```

Имя модуля — обычно путь к репозиторию. Это важно для импорта пакетов из модуля другими.

---

## Файл go.mod

```go
module github.com/myuser/myapp

go 1.26

require (
    github.com/gin-gonic/gin v1.9.1
    go.uber.org/zap v1.27.0
    golang.org/x/text v0.14.0  // indirect
)
```

### Директивы go.mod

**`module`** — путь модуля (один, обязательный):
```go
module github.com/mycompany/myservice
```

**`go`** — минимальная версия Go:
```go
go 1.26
```

**`require`** — зависимости с минимальной версией:
```go
require (
    github.com/pkg/errors v0.9.1
    golang.org/x/crypto v0.21.0 // indirect
)
```
`// indirect` — транзитивная зависимость (нужна зависимости вашей зависимости).

**`replace`** — замена пакета (локальная разработка или форк):
```go
replace (
    github.com/original/pkg => ./local/pkg   // локальная замена
    github.com/original/pkg v1.2.3 => github.com/fork/pkg v1.2.4  // форк
)
```

**`exclude`** — исключение конкретной версии:
```go
exclude github.com/broken/pkg v1.3.0  // эта версия содержит баг
```

---

## Файл go.sum

`go.sum` — файл с криптографическими хешами всех зависимостей. Обеспечивает воспроизводимые сборки:

```
github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Tnip5kqbEAP+KLuI3SUcPTeU=
```

**Не редактируй go.sum вручную** — он обновляется командами `go` автоматически. Коммить go.sum в репозиторий — обязательно.

---

## Основные команды

### Добавление зависимости

```bash
go get github.com/pkg/errors           # последняя версия
go get github.com/pkg/errors@v0.9.1    # конкретная версия
go get github.com/pkg/errors@latest    # явно последняя
go get github.com/pkg/errors@v0.9.0+incompatible  # pre-modules версия
```

После `go get` зависимость добавляется в `go.mod` и `go.sum`.

### go mod tidy — синхронизация

```bash
go mod tidy
```

Самая важная команда для поддержания порядка:
- Добавляет пропущенные зависимости
- Удаляет неиспользуемые зависимости
- Обновляет `go.sum`

**Запускай после каждого изменения импортов.**

### go mod vendor — локальная копия зависимостей

```bash
go mod vendor
```

Копирует все зависимости в директорию `vendor/`. Полезно для:
- Окружений без интернета (CI/CD)
- Гарантии, что зависимости не исчезнут с хостинга
- Одобрения зависимостей security-командой

```bash
# Сборка из vendor:
go build -mod=vendor ./...

# Тесты из vendor:
go test -mod=vendor ./...
```

---

## Семантическое версионирование

Go Modules следуют semver: `vMAJOR.MINOR.PATCH`

```
v1.2.3
│ │ └── Patch: исправление ошибок (обратно совместимо)
│ └──── Minor: новые фичи (обратно совместимо)
└────── Major: breaking changes (несовместимо)
```

### Важное правило: major version ≥ 2

Если модуль достиг v2+, его путь должен содержать major version:

```
v1: github.com/foo/bar
v2: github.com/foo/bar/v2
v3: github.com/foo/bar/v3
```

В коде:
```go
import (
    bar "github.com/foo/bar"       // v1
    barv2 "github.com/foo/bar/v2"  // v2 — разные пути, можно использовать оба!
)
```

---

## Minimal Version Selection (MVS)

Go использует MVS — детерминированный алгоритм выбора версий:
- Выбирает минимальную версию, удовлетворяющую всем требованиям
- Предсказуемый результат без «случайных» обновлений

```bash
go list -m all  # посмотреть все версии в сборке
```

---

## Практический пример: multi-module workspace

Когда разрабатываешь несколько связанных модулей локально — используй Go Workspaces (Go 1.18+):

```bash
mkdir workspace && cd workspace
mkdir mod1 mod2

cd mod1 && go mod init example.com/mod1
cd ../mod2 && go mod init example.com/mod2

cd .. && go work init mod1 mod2
```

Создаётся `go.work`:
```go
go 1.26

use (
    ./mod1
    ./mod2
)
```

Теперь `mod2` может импортировать `mod1` без `replace` директив.

---

## Итог

- `go mod init path` — создать новый модуль
- `go.mod` — зависимости и версии; `go.sum` — хеши для integrity
- `go mod tidy` — синхронизировать go.mod с реальными импортами
- `go mod vendor` — скопировать зависимости в vendor/
- `go get pkg@version` — добавить или обновить зависимость
- Semver: v1 совместим, v2+ меняет путь модуля
- `replace` — для локальной разработки и форков
