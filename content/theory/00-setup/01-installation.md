---
title: "Установка Go"
description: "Пошаговая установка Go на Windows, macOS и Linux, настройка окружения и первая программа"
order: 1
---

# Установка Go

Прежде чем писать код, нужно один раз правильно настроить окружение. Это займёт 10–15 минут и избавит от множества головных болей в будущем.

## Где скачать

Официальная страница загрузки: **https://go.dev/dl/**

Актуальная стабильная версия на момент написания: **Go 1.26.2**.

Всегда скачивай с официального сайта — не используй пакетные менеджеры вроде apt/yum/brew для production-разработки, там часто устаревшие версии.

---

## Установка на Linux

### Шаг 1. Скачай архив

```bash
wget https://go.dev/dl/go1.26.2.linux-amd64.tar.gz
```

Если у тебя ARM-процессор (например, Raspberry Pi или Apple Silicon через Rosetta), замени `amd64` на `arm64`.

### Шаг 2. Удали старую версию и распакуй новую

```bash
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.2.linux-amd64.tar.gz
```

> **Важно**: не распаковывай в уже существующую директорию `/usr/local/go` — всегда сначала удаляй старую, иначе получишь смесь файлов от разных версий.

### Шаг 3. Добавь Go в PATH

Открой файл `~/.profile` (или `~/.bashrc` / `~/.zshrc` в зависимости от оболочки) и добавь в конец:

```bash
export PATH=$PATH:/usr/local/go/bin
```

Применить изменения без перезапуска терминала:

```bash
source ~/.profile
```

### Шаг 4. Проверь установку

```bash
go version
# go version go1.26.2 linux/amd64
```

---

## Установка на macOS

### Вариант 1: PKG-пакет (рекомендуется)

1. Скачай файл `go1.26.2.darwin-amd64.pkg` (или `darwin-arm64.pkg` для Apple Silicon) с https://go.dev/dl/
2. Открой скачанный `.pkg` файл
3. Следуй инструкциям установщика — Go установится в `/usr/local/go`
4. Установщик **автоматически** добавит `/usr/local/go/bin` в PATH через `/etc/paths.d/go`
5. Перезапусти терминал

### Вариант 2: Вручную через архив (аналогично Linux)

```bash
# Для Apple Silicon (M1/M2/M3/M4)
curl -Lo go1.26.2.darwin-arm64.tar.gz https://go.dev/dl/go1.26.2.darwin-arm64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.2.darwin-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
source ~/.zshrc
```

### Проверь установку

```bash
go version
# go version go1.26.2 darwin/arm64
```

---

## Установка на Windows

### Шаг 1. Скачай MSI-установщик

Скачай файл `go1.26.2.windows-amd64.msi` с https://go.dev/dl/

### Шаг 2. Запусти установщик

1. Дважды кликни на скачанный `.msi` файл
2. Следуй инструкциям — Go установится в `C:\Program Files\Go\`
3. Установщик **автоматически** добавит `C:\Program Files\Go\bin` в переменную окружения `PATH`

### Шаг 3. Перезапусти командную строку

Закрой и снова открой `cmd` или PowerShell — иначе новый PATH не применится.

### Шаг 4. Проверь установку

```cmd
go version
:: go version go1.26.2 windows/amd64
```

---

## GOPATH и GOROOT — когда они нужны

Один из частых источников путаницы для новичков.

### GOROOT

Это директория, куда установлен сам Go (`/usr/local/go`). **Не нужно устанавливать вручную** — Go находит себя автоматически. Не трогай эту переменную, если только у тебя нет нескольких версий Go.

### GOPATH

Исторически это была рабочая директория для Go-проектов. **Сейчас (с появлением Go Modules в Go 1.11) GOPATH практически не нужен** для повседневной разработки.

Если не задан явно, по умолчанию:
- Linux/macOS: `~/go`
- Windows: `%USERPROFILE%\go`

**Когда GOPATH всё ещё важен:**
- Утилиты, установленные через `go install`, попадают в `$GOPATH/bin` — добавь этот путь в PATH:

```bash
# Linux/macOS — добавь в ~/.profile или ~/.zshrc
export PATH=$PATH:$(go env GOPATH)/bin
```

```cmd
:: Windows PowerShell
$env:PATH += ";$(go env GOPATH)\bin"
```

Проверить текущие значения:

```bash
go env GOPATH
go env GOROOT
```

---

## Настройка VS Code

VS Code — наиболее популярный редактор для Go среди новичков и опытных разработчиков.

### Установка расширения

1. Установи VS Code: https://code.visualstudio.com/
2. Установи официальное расширение Go: https://marketplace.visualstudio.com/items?itemName=golang.go
   Или через терминал: `code --install-extension golang.go`

### Установка инструментов Go

При первом открытии `.go` файла VS Code предложит установить дополнительные инструменты (gopls, dlv, staticcheck и т.д.). Нажми **"Install All"**.

Либо вручную:

```bash
go install golang.org/x/tools/gopls@latest
```

### Рекомендуемые настройки VS Code

Добавь в `settings.json` (Ctrl+Shift+P → "Open Settings JSON"):

```json
{
  "editor.formatOnSave": true,
  "go.useLanguageServer": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  }
}
```

---

## Первая программа: Hello, World!

### Шаг 1. Создай рабочую папку

```bash
mkdir ~/projects/hello
cd ~/projects/hello
```

### Шаг 2. Инициализируй модуль

```bash
go mod init hello
```

Это создаст файл `go.mod` — основу системы модулей Go. Подробнее разберём в главе 9 «Пакеты и модули».

### Шаг 3. Создай файл main.go

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Шаг 4. Запусти

```bash
go run main.go
# Hello, World!
```

---

## go run vs go build vs go install

Три команды, которые часто путают.

### go run

```bash
go run main.go
```

Компилирует и сразу запускает — временный бинарник удаляется после выполнения. Удобно для быстрых экспериментов и обучения.

### go build

```bash
go build -o hello main.go
# Создаёт бинарный файл ./hello (или hello.exe на Windows)
./hello
# Hello, World!
```

Компилирует в постоянный исполняемый файл. Используй для сборки программы для деплоя или распространения.

### go install

```bash
go install .
# Компилирует и кладёт бинарник в $GOPATH/bin/hello
```

Аналогично `go build`, но кладёт результат в `$GOPATH/bin`. Используется для установки утилит (например, `go install golang.org/x/tools/cmd/goimports@latest`).

---

## Проверка окружения

Полезная команда для диагностики:

```bash
go env
```

Выводит все переменные окружения Go. Если что-то пошло не так с настройкой PATH или GOPATH — проверяй здесь.

```bash
go env GOPATH     # Рабочая директория Go
go env GOROOT     # Директория установки Go
go env GOOS       # Целевая ОС (linux, darwin, windows)
go env GOARCH     # Целевая архитектура (amd64, arm64)
```

---

## Типичные ошибки при установке

**`go: command not found`** — PATH не настроен. Проверь, что `/usr/local/go/bin` есть в `$PATH` и что ты перезапустил терминал после изменения `.profile`.

**Старая версия после обновления** — Возможно, в системе несколько установок Go. Проверь: `which go` (Linux/macOS) или `where go` (Windows), удали старую установку.

**Ошибки прав доступа на Linux/macOS** — При распаковке архива используй `sudo`. Файлы в `/usr/local/go` должны принадлежать root.

---

## Итог

Ты установил Go, настроил редактор и запустил первую программу. Теперь окружение готово к работе.

**Что запомнить:**
- Скачивай только с официального `go.dev/dl/`
- GOPATH нужен только для `go install` — добавь `$GOPATH/bin` в PATH
- `go run` — для экспериментов, `go build` — для production, `go install` — для утилит
- При любых проблемах — `go env` покажет текущее состояние окружения
