---
title: "Модели данных"
difficulty: easy
order: 1
file: "model/task.go"
hints:
  - "Используй struct с JSON тегами для каждого поля"
  - "Status лучше сделать отдельным типом: type Status string"
  - "Конструктор NewTask должен проставлять StatusPending и CreatedAt автоматически"
---

# Модели данных

В Go приложениях **модели** — это структуры данных, описывающие сущности вашего домена. Они используются всеми слоями: базой данных, бизнес-логикой и HTTP handlers.

Хорошо спроектированная модель:
- Описывает структуру данных через поля с правильными типами
- Содержит JSON-теги для сериализации/десериализации
- Инкапсулирует правила создания через конструкторы
- Валидирует данные перед использованием

## Что нужно сделать

Создайте модель задачи в пакете `model`.

### 1. Тип `Status`

Определи именованный тип `Status` на основе `string` с тремя константами:
- `StatusPending = "pending"`
- `StatusInProgress = "in_progress"`
- `StatusDone = "done"`

Именованный тип (не просто строка) позволяет компилятору проверять корректность статусов на этапе компиляции.

### 2. Структура `Task`

| Поле        | Тип         | JSON-тег         |
|-------------|-------------|-----------------|
| ID          | `int64`     | `"id"`          |
| Title       | `string`    | `"title"`       |
| Description | `string`    | `"description"` |
| Status      | `Status`    | `"status"`      |
| CreatedAt   | `time.Time` | `"created_at"`  |
| UpdatedAt   | `time.Time` | `"updated_at"`  |

### 3. Конструктор `NewTask`

```go
func NewTask(title, description string) Task
```

Автоматически устанавливает `Status = StatusPending` и заполняет `CreatedAt`, `UpdatedAt` текущим временем через `time.Now()`.

### 4. Метод `Validate`

```go
func (t Task) Validate() error
```

Возвращает `errors.New("title is required")` если `Title` пустой, иначе `nil`.

## Требования

- Все поля структуры должны иметь JSON-теги в snake_case
- `NewTask` должен всегда устанавливать `StatusPending` — клиент не может задать начальный статус
- `Validate` должен возвращать именно `errors.New("title is required")` при пустом Title

## Пример использования

```go
task := model.NewTask("Написать тесты", "Покрыть все хендлеры")
// task.Status == model.StatusPending
// task.CreatedAt != time.Time{}

if err := task.Validate(); err != nil {
    log.Fatal(err)
}

// Пустой title не пройдёт валидацию
bad := model.Task{}
fmt.Println(bad.Validate()) // title is required
```
