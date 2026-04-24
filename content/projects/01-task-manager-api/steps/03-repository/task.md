---
title: "Repository — слой работы с PostgreSQL"
difficulty: medium
order: 3
file: "repository/task_repo.go"
hints:
  - "В PostgreSQL плейсхолдеры — $1, $2, $3, а не ? как в SQLite"
  - "Для INSERT используй RETURNING id, created_at, updated_at — получишь значения за один запрос"
  - "При GetByID проверяй sql.ErrNoRows и оборачивай в понятную ошибку через fmt.Errorf"
  - "rows.Err() нужно проверять после цикла for rows.Next()"
---

# Repository — слой работы с PostgreSQL

**Repository** — паттерн, который инкапсулирует всю логику доступа к данным. Бизнес-логика не знает, откуда приходят данные: PostgreSQL, Redis, файл или in-memory. Она просто вызывает методы репозитория.

Преимущества:
- Бизнес-логика изолирована от деталей хранения
- Легко тестировать (заменяй реальный репозиторий фейком)
- Изменение СУБД не затрагивает остальной код

## Что нужно сделать

Создай `TaskRepository` с пятью методами. Все методы принимают `context.Context` первым аргументом — это идиоматичный Go для операций с I/O.

### Структура и конструктор

```go
type TaskRepository struct {
    db *sql.DB
}

func New(db *sql.DB) *TaskRepository
```

### Методы

**`Create`** — вставляет задачу, получает ID и времена из базы:
```go
func (r *TaskRepository) Create(ctx context.Context, task model.Task) (model.Task, error)
```
Используй `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

**`GetByID`** — находит задачу по ID:
```go
func (r *TaskRepository) GetByID(ctx context.Context, id int64) (model.Task, error)
```
При `sql.ErrNoRows` верни `fmt.Errorf("task %d: not found", id)`

**`List`** — все задачи, новые сначала:
```go
func (r *TaskRepository) List(ctx context.Context) ([]model.Task, error)
```
Используй `ORDER BY created_at DESC`. Возвращай пустой slice `make([]model.Task, 0)`, не nil.

**`Update`** — обновляет задачу, получает новое `updated_at`:
```go
func (r *TaskRepository) Update(ctx context.Context, task model.Task) (model.Task, error)
```
Используй `UPDATE ... SET ..., updated_at=NOW() WHERE id=$4 RETURNING updated_at`

**`Delete`** — удаляет задачу:
```go
func (r *TaskRepository) Delete(ctx context.Context, id int64) error
```
Проверяй `result.RowsAffected()` — если 0, задачи не существовало.

## Требования

- Все плейсхолдеры в PostgreSQL: `$1`, `$2`, не `?`
- `Create` должен получать `id`, `created_at`, `updated_at` через `RETURNING`
- `List` возвращает `[]model.Task{}`, не `nil`
- Проверяй `rows.Err()` после цикла `for rows.Next()`

## Пример использования

```go
repo := repository.New(db)

task := model.NewTask("Deploy app", "Настроить CI/CD")
created, err := repo.Create(ctx, task)
// created.ID заполнен из базы

tasks, err := repo.List(ctx)
// []model.Task с записями из БД
```
