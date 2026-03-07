---
title: "Сервисный слой"
description: "Реализуйте бизнес-логику для управления задачами"
order: 3
difficulty: medium
file: "service/todo.go"
hints:
  - "Сервис не работает с данными напрямую — он вызывает методы storage"
  - "В Create создайте model.Todo из CreateTodoRequest и передайте в storage.Create"
  - "В Update сначала получите существующую задачу через GetByID, затем обновите только переданные поля (проверяйте указатели на nil)"
  - "Не забудьте обработать ошибки от storage — если задача не найдена, пробросьте ошибку выше"
---

# Сервисный слой (Service)

**Сервис** — это слой бизнес-логики. Он стоит между HTTP-обработчиками и хранилищем. Сервис:
- Преобразует данные запроса в модели
- Содержит бизнес-правила
- Вызывает методы хранилища

## Зачем нужен отдельный сервис?

Без сервиса обработчики HTTP напрямую работают с хранилищем. Это приводит к проблемам:
- Бизнес-логика размазана по обработчикам
- Сложно писать unit-тесты
- Нельзя переиспользовать логику

## Что нужно реализовать

### Структура `TodoService`

```go
type TodoService struct {
    storage *storage.TodoStorage
}
```

### Конструктор

```go
func NewTodoService(s *storage.TodoStorage) *TodoService
```

### Методы

| Метод | Сигнатура | Описание |
|-------|-----------|----------|
| Create | `Create(req model.CreateTodoRequest) model.Todo` | Создаёт Todo из запроса, передаёт в storage |
| GetByID | `GetByID(id int) (model.Todo, error)` | Делегирует в storage |
| GetAll | `GetAll() []model.Todo` | Делегирует в storage |
| Update | `Update(id int, req model.UpdateTodoRequest) (model.Todo, error)` | Частичное обновление |
| Delete | `Delete(id int) error` | Делегирует в storage |

### Логика Update

Метод `Update` — ключевой. Он реализует **частичное обновление**:

```go
func (s *TodoService) Update(id int, req model.UpdateTodoRequest) (model.Todo, error) {
    // 1. Получить существующую задачу
    // 2. Обновить только поля, где указатель != nil
    // 3. Сохранить через storage.Update
}
```

Пример: если передан `{"done": true}`, то Title и Description не меняются.
