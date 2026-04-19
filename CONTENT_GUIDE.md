# Руководство по структуре контента Go Path

Это руководство объясняет, как заполняется теория, quiz и другой контент в проекте Go Path для Claude Code и других разработчиков.

---

## 📚 Структура теории (Theory)

### Расположение файлов
```
content/theory/
├── 01-basics/
│   ├── meta.yaml
│   ├── quiz.yaml
│   ├── 01-hello-world.md
│   ├── 02-variables.md
│   └── ...
├── 02-types/
│   ├── meta.yaml
│   ├── quiz.yaml
│   └── ...
└── 03-functions/
    ├── meta.yaml
    ├── quiz.yaml
    └── ...
```

### 1. Файл meta.yaml (метаданные главы)

Каждая папка главы должна содержать `meta.yaml` с описанием главы:

```yaml
title: "Основы Go"
description: "Введение в язык Go: первая программа, переменные, константы"
order: 1
```

**Поля:**
- `title` (обязательно) — название главы, отображается в UI
- `description` (обязательно) — краткое описание главы
- `order` (обязательно) — порядковый номер для сортировки глав

### 2. Файлы уроков (*.md)

Каждый урок — это отдельный markdown файл с frontmatter.

**Структура файла урока:**

```markdown
---
title: "Hello, World!"
description: "Первая программа на Go — структура файла, пакет main и функция main"
order: 1
---

# Hello, World!

Каждая программа на Go начинается с объявления пакета...

## Минимальная программа

\`\`\`go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
\`\`\`

## Разбор

- **`package main`** — указывает, что это исполняемый пакет
- **`import "fmt"`** — подключает стандартный пакет
```

**Frontmatter поля:**
- `title` (обязательно) — название урока
- `description` (обязательно) — краткое описание урока
- `order` (обязательно) — порядковый номер урока в главе

**Контент:**
- После frontmatter (после второго `---`) идёт markdown контент
- Используйте заголовки, код-блоки, списки, цитаты
- Код-блоки должны указывать язык: \`\`\`go

**Именование файлов:**
- Формат: `NN-slug-name.md` (например, `01-hello-world.md`)
- Slug используется в URL и API
- Номер в имени файла помогает с сортировкой в файловой системе

### 3. Как работает загрузка теории

**Код загрузки** (`internal/service/theory.go`):

1. **Сканирование папок** — система читает все папки в `content/theory/`
2. **Чтение meta.yaml** — загружает метаданные главы
3. **Чтение .md файлов** — парсит frontmatter и контент каждого урока
4. **Сортировка** — главы и уроки сортируются по полю `order`
5. **Создание индексов** — создаются карты для быстрого доступа по slug

**Структура данных:**
```go
type Chapter struct {
    Slug        string   // имя папки (например, "01-basics")
    Title       string   // из meta.yaml
    Description string   // из meta.yaml
    Order       int      // из meta.yaml
    Lessons     []Lesson // все уроки главы
}

type Lesson struct {
    Slug        string // имя файла без .md (например, "01-hello-world")
    Title       string // из frontmatter
    Description string // из frontmatter
    Order       int    // из frontmatter
    ChapterSlug string // ссылка на главу
    Content     string // markdown контент после frontmatter
}
```

### 4. API endpoints для теории

- `GET /api/theory/chapters` — список всех глав с уроками (без контента)
- `GET /api/theory/chapters/:slug` — одна глава с уроками (без контента)
- `GET /api/theory/chapters/:chapterSlug/lessons/:lessonSlug` — один урок с полным контентом
- `POST /api/theory/chapters/:chapterSlug/lessons/:lessonSlug/complete` — отметить урок как прочитанный

---

## 🎯 Структура Quiz

### Расположение файлов

Quiz находятся в тех же папках, что и теория:

```
content/theory/
├── 01-basics/
│   ├── meta.yaml
│   ├── quiz.yaml          ← Quiz для главы
│   ├── 01-hello-world.md
│   └── ...
```

### Файл quiz.yaml

**Структура:**

```yaml
questions:
  - question: "Какой функцией выводят текст в консоль?"
    options:
      - "fmt.Println"
      - "console.log"
      - "print()"
      - "System.out.println"
    answer: 0
    explanation: "В Go для вывода используется пакет fmt. Функция fmt.Println выводит строку с переводом строки."

  - question: "Какое имя пакета должно быть у исполняемой программы на Go?"
    options:
      - "app"
      - "main"
      - "program"
      - "run"
    answer: 1
    explanation: "Исполняемые программы в Go всегда принадлежат пакету main."

  - question: "Какой командой запускается Go-программа?"
    options:
      - "go start main.go"
      - "go exec main.go"
      - "go run main.go"
      - "go build main.go && ./main"
    answer: 2
    explanation: "Команда go run компилирует и сразу запускает программу."
```

**Поля вопроса:**
- `question` (обязательно) — текст вопроса
- `options` (обязательно) — массив из 4 вариантов ответа
- `answer` (обязательно) — индекс правильного ответа (0-3)
- `explanation` (обязательно) — объяснение правильного ответа

**Важные правила:**
1. **Индексация с 0** — первый вариант имеет индекс 0, второй — 1, и т.д.
2. **Всегда 4 варианта** — для единообразия UI
3. **Комментарии** — можно добавлять комментарии для группировки вопросов по урокам

### Как работает загрузка Quiz

**Код загрузки** (`internal/service/quiz.go`):

1. **Сканирование папок** — читает все папки в `content/theory/`
2. **Чтение meta.yaml** — получает название главы
3. **Чтение quiz.yaml** — парсит вопросы
4. **Генерация ID** — каждому вопросу присваивается уникальный ID: `{chapterSlug}:{index}`
5. **Создание индексов** — вопросы индексируются по главам и по ID

**Структура данных:**
```go
type QuizQuestion struct {
    ID          string   // генерируется: "01-basics:0"
    Question    string   // текст вопроса
    Options     []string // варианты ответа
    Answer      int      // индекс правильного ответа (НЕ отправляется клиенту)
    Explanation string   // объяснение (НЕ отправляется до проверки)
    ChapterSlug string   // ссылка на главу
}

type QuizChapterInfo struct {
    Slug          string // slug главы
    Title         string // название главы
    QuestionCount int    // количество вопросов
}
```

### API endpoints для Quiz

- `GET /api/quiz/chapters` — список глав с количеством вопросов
- `GET /api/quiz?chapters=01-basics,02-types&limit=10` — получить случайные вопросы
- `POST /api/quiz/answer` — проверить ответ
  ```json
  {
    "question_id": "01-basics:0",
    "answer": 0
  }
  ```
  Ответ:
  ```json
  {
    "correct": true,
    "correct_answer": 0,
    "explanation": "В Go для вывода используется пакет fmt..."
  }
  ```

### Рекомендации по созданию Quiz

1. **Связь с уроками** — вопросы должны покрывать материал из уроков главы
2. **Разнообразие** — смешивайте типы вопросов:
   - Определения и концепции
   - Синтаксис и код
   - Поведение и результаты
   - Best practices
3. **Качество вариантов** — неправильные варианты должны быть правдоподобными
4. **Объяснения** — всегда объясняйте, почему ответ правильный
5. **Количество** — рекомендуется 10-20 вопросов на главу

---

## 📝 Структура задач (Tasks)

### Расположение файлов

```
content/tasks/
├── 01-basics/
│   ├── meta.yaml
│   ├── 01-hello-name/
│   │   ├── task.md
│   │   ├── template.go
│   │   ├── solution_test.go
│   │   └── completions.yaml
│   └── 02-even-odd/
│       └── ...
```

### Файлы задачи

**1. task.md** — описание задачи с frontmatter:

```markdown
---
title: "Hello, Name"
difficulty: "easy"
order: 1
---

Напишите функцию `greet`, которая принимает имя и возвращает приветствие.

## Примеры

\`\`\`go
greet("Gopher") // "Hello, Gopher!"
greet("World")  // "Hello, World!"
\`\`\`
```

**2. template.go** — начальный код для пользователя:

```go
package main

// greet возвращает приветствие для указанного имени
func greet(name string) string {
	// TODO: реализуйте функцию
	return ""
}
```

**3. solution_test.go** — тесты для проверки:

```go
package main

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Gopher", "Hello, Gopher!"},
		{"World", "Hello, World!"},
	}
	
	for _, tt := range tests {
		got := greet(tt.name)
		if got != tt.want {
			t.Errorf("greet(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}
```

**4. completions.yaml** — подсказки для автодополнения:

```yaml
completions:
  - trigger: "greet"
    code: |
      func greet(name string) string {
      	return "Hello, " + name + "!"
      }
    description: "Полное решение функции greet"
```

---

## 🚀 Структура проектов (Projects)

### Расположение файлов

```
content/projects/
├── 01-todo-api/
│   ├── meta.yaml
│   ├── go.mod.tmpl
│   └── steps/
│       ├── 01-models/
│       │   ├── task.md
│       │   ├── template.go
│       │   ├── completions.yaml
│       │   ├── reference/
│       │   │   └── model/
│       │   │       └── todo.go
│       │   └── tests/
│       │       └── model/
│       │           └── todo_test.go
│       └── 02-storage/
│           └── ...
```

### Файлы проекта

**1. meta.yaml** — метаданные проекта:

```yaml
title: "TODO API"
description: "Создание REST API для управления задачами"
order: 1
```

**2. go.mod.tmpl** — шаблон go.mod для проекта:

```
module todoapi

go 1.21

require (
	github.com/gorilla/mux v1.8.0
)
```

**3. Шаги проекта** — каждый шаг в папке `steps/NN-name/`:

**task.md** с frontmatter:
```markdown
---
title: "Модели данных"
difficulty: "easy"
order: 1
file: "model/todo.go"
hints:
  - "Используйте struct для определения модели"
  - "Добавьте JSON теги для сериализации"
---

Создайте модель данных для TODO задачи...
```

**template.go** — начальный код
**reference/** — эталонные файлы для следующих шагов
**tests/** — тесты для проверки

---

## 🔧 Как добавить новый контент

### Добавить новый урок теории

1. Перейдите в нужную главу: `content/theory/NN-chapter/`
2. Создайте файл `NN-lesson-name.md`
3. Добавьте frontmatter с полями `title`, `description`, `order`
4. Напишите markdown контент
5. Обновите `quiz.yaml` — добавьте вопросы по новому уроку

### Добавить новую главу теории

1. Создайте папку `content/theory/NN-new-chapter/`
2. Создайте `meta.yaml` с метаданными главы
3. Создайте `quiz.yaml` с вопросами
4. Добавьте файлы уроков `*.md`
5. Перезапустите сервер — контент загрузится автоматически

### Добавить новую задачу

1. Перейдите в главу: `content/tasks/NN-chapter/`
2. Создайте папку `NN-task-name/`
3. Создайте файлы:
   - `task.md` — описание
   - `template.go` — начальный код
   - `solution_test.go` — тесты
   - `completions.yaml` — автодополнения
4. Перезапустите сервер

### Добавить новый проект

1. Создайте папку `content/projects/NN-project-name/`
2. Создайте `meta.yaml` и `go.mod.tmpl`
3. Создайте папку `steps/`
4. Для каждого шага создайте папку со структурой:
   - `task.md`, `template.go`, `completions.yaml`
   - `reference/` — эталонные файлы
   - `tests/` — тесты
5. Перезапустите сервер

---

## 📊 Диагностика и отладка

### Проверка загрузки контента

При запуске сервера в логах появляются сообщения:

```
INFO theory loaded {"chapters": 12, "total_lessons": 45}
INFO quiz loaded {"chapters": 12, "total_questions": 180}
INFO tasks loaded {"chapters": 5, "total_tasks": 25}
INFO projects loaded {"count": 2, "total_steps": 8}
```

### Частые ошибки

1. **Отсутствует frontmatter** — файл должен начинаться с `---`
2. **Неверный YAML** — проверьте отступы и синтаксис
3. **Отсутствует meta.yaml** — глава/проект будет пропущена
4. **Неверный индекс answer** — помните, что индексация с 0
5. **Отсутствуют обязательные файлы** — task.md, template.go, solution_test.go

### Логи предупреждений

Если файл не загружается, в логах появится:

```
WARN skipping chapter {"dir": "03-broken", "error": "missing meta.yaml"}
WARN skipping lesson {"file": "01-bad.md", "error": "invalid frontmatter"}
WARN skipping quiz {"chapter": "01-basics", "error": "yaml: unmarshal error"}
```

---

## 🎨 Best Practices

### Теория

1. **Структура** — начинайте с простого, постепенно усложняйте
2. **Примеры** — каждая концепция должна иметь код-пример
3. **Практичность** — показывайте реальные use cases
4. **Советы** — используйте блоки цитат для важных замечаний

### Quiz

1. **Покрытие** — вопросы должны покрывать весь материал главы
2. **Баланс** — смешивайте простые и сложные вопросы
3. **Актуальность** — вопросы должны быть связаны с уроками
4. **Объяснения** — всегда объясняйте правильный ответ

### Задачи

1. **Градация** — от простых к сложным
2. **Тесты** — покрывайте edge cases
3. **Шаблон** — давайте достаточно кода для старта
4. **Описание** — четко формулируйте требования

### Проекты

1. **Инкрементальность** — каждый шаг добавляет новую функциональность
2. **Независимость** — шаги должны быть относительно независимыми
3. **Реалистичность** — проект должен быть похож на реальный
4. **Документация** — объясняйте архитектурные решения

---

## 🔄 Автоматическая перезагрузка

Контент загружается при старте сервера. Для применения изменений:

1. Остановите сервер (Ctrl+C)
2. Внесите изменения в файлы
3. Запустите сервер снова: `go run cmd/server/main.go`

В будущем можно добавить hot reload для автоматической перезагрузки.

---

## 📚 Дополнительные ресурсы

- **Код загрузки теории**: `internal/service/theory.go`
- **Код загрузки quiz**: `internal/service/quiz.go`
- **Код загрузки задач**: `internal/service/task.go`
- **Код загрузки проектов**: `internal/service/project.go`
- **Модели данных**: `internal/model/`

---

## ✅ Чеклист для Claude Code

При создании нового контента проверьте:

- [ ] Frontmatter присутствует и корректен
- [ ] Все обязательные поля заполнены
- [ ] YAML синтаксис валиден
- [ ] Порядковые номера (order) уникальны и последовательны
- [ ] Код-примеры синтаксически корректны
- [ ] Тесты проходят
- [ ] Quiz вопросы имеют правильные индексы ответов (0-3)
- [ ] Объяснения в quiz понятны и полны
- [ ] Файловая структура соответствует шаблону
- [ ] Нет опечаток и грамматических ошибок

---

**Версия документа**: 1.0  
**Дата обновления**: 2026-04-19
