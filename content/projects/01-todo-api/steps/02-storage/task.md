---
title: "In-Memory хранилище"
description: "Реализуйте потокобезопасное хранилище задач в памяти"
order: 2
difficulty: medium
file: "storage/todo.go"
hints:
  - "Используйте sync.RWMutex для потокобезопасности: RLock для чтения, Lock для записи"
  - "Храните задачи в map[int]model.Todo с автоинкрементом ID"
  - "В Create устанавливайте ID (nextID++) и CreatedAt (time.Now())"
  - "GetAll должен возвращать пустой слайс (не nil), если задач нет"
  - "Не забудьте определить переменную ошибки: var ErrNotFound = errors.New(\"todo not found\")"
---

# In-Memory хранилище (Storage)

**Хранилище (Storage/Repository)** — это слой, отвечающий за сохранение и извлечение данных. В реальных приложениях здесь работают с базой данных, но для начала мы реализуем хранилище в памяти.

## Паттерн Repository

Идея проста: весь код, работающий с данными, изолирован в одном месте. Сервисный слой не знает, где хранятся данные — в памяти, PostgreSQL или файле.

## Что нужно реализовать

### Ошибка `ErrNotFound`

```go
var ErrNotFound = errors.New("todo not found")
```

### Структура `TodoStorage`

Хранилище с потокобезопасным доступом:

```go
type TodoStorage struct {
    mu     sync.RWMutex
    todos  map[int]model.Todo
    nextID int
}
```

### Конструктор `NewTodoStorage`

Возвращает готовое к использованию хранилище с инициализированной картой и `nextID = 1`.

### Методы

| Метод | Сигнатура | Описание |
|-------|-----------|----------|
| Create | `Create(todo model.Todo) model.Todo` | Присваивает ID, устанавливает CreatedAt, сохраняет |
| GetByID | `GetByID(id int) (model.Todo, error)` | Возвращает задачу или ErrNotFound |
| GetAll | `GetAll() []model.Todo` | Все задачи (пустой слайс, если нет) |
| Update | `Update(id int, todo model.Todo) (model.Todo, error)` | Обновляет или ErrNotFound |
| Delete | `Delete(id int) error` | Удаляет или ErrNotFound |

## Потокобезопасность

Используйте `sync.RWMutex`:
- **RLock/RUnlock** — для операций чтения (GetByID, GetAll)
- **Lock/Unlock** — для операций записи (Create, Update, Delete)

```go
func (s *TodoStorage) GetByID(id int) (model.Todo, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    // ...
}
```
