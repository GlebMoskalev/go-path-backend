---
title: "Интеграционные тесты"
difficulty: hard
order: 7
file: "integration_test.go"
hints:
  - "testcontainers поднимает реальный PostgreSQL в Docker — никаких моков"
  - "httptest.NewServer запускает настоящий HTTP сервер для тестов"
  - "Каждый тест получает свежую БД через setupTestServer"
  - "t.Cleanup регистрирует завершение контейнера и сервера автоматически"
---

# Интеграционные тесты

Интеграционные тесты проверяют весь стек приложения как единое целое: реальная БД, реальный HTTP сервер, реальные запросы. Это защита от ошибок на стыке компонентов.

В отличие от unit-тестов с моками, интеграционные тесты:
- Находят ошибки несовместимости между слоями
- Тестируют SQL-запросы на реальной СУБД
- Проверяют сериализацию/десериализацию JSON end-to-end

`testcontainers-go` запускает настоящий PostgreSQL в Docker-контейнере. Контейнер создаётся для каждого теста и автоматически завершается через `t.Cleanup`.

## Что нужно реализовать

### Вспомогательная функция `setupTestServer`

```go
func setupTestServer(t *testing.T) *httptest.Server
```

Поднимает полный стек:
1. PostgreSQL контейнер через testcontainers
2. `db.New` + `db.Migrate`
3. `repository.New`
4. `handler.New`
5. `server.NewRouter` обёрнутый в `middleware.Chain(router, middleware.Logger, middleware.Recovery)`
6. `httptest.NewServer(handler)` — настоящий HTTP сервер на случайном порту

Регистрируй завершение через `t.Cleanup`.

### Тесты (минимум 6)

| Тест | Что проверяет |
|------|---------------|
| `TestHealthCheck` | GET /health → 200, body `{"status":"ok"}` |
| `TestCreateAndGetTask` | POST создаёт задачу, GET /tasks/{id} возвращает её |
| `TestListTasks` | Создать 3 задачи, GET /tasks возвращает все |
| `TestUpdateTask` | Создать, PUT обновить, проверить изменения |
| `TestDeleteTask` | Создать, DELETE, GET возвращает 404 |
| `TestValidation` | POST с пустым title → 400 |

## Требования

- Каждый тест вызывает `setupTestServer(t)` — изолированная среда
- Используй `http.DefaultClient` для HTTP-запросов к тестовому серверу
- Для `POST/PUT` устанавливай `Content-Type: application/json`
- Проверяй и статус коды, и содержимое тела ответа

## Пример структуры

```go
func TestCreateAndGetTask(t *testing.T) {
    srv := setupTestServer(t)

    // создаём задачу
    resp, err := http.Post(srv.URL+"/tasks", "application/json",
        strings.NewReader(`{"title":"Test","description":"Desc"}`))
    // ...проверяем resp.StatusCode == 201

    // получаем по ID
    var created model.Task
    json.NewDecoder(resp.Body).Decode(&created)
    resp2, _ := http.Get(fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID))
    // ...проверяем resp2.StatusCode == 200
}
```
