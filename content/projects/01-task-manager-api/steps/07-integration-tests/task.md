---
title: "Интеграционные тесты"
difficulty: hard
order: 7
file: "integration/setup_test.go"
hints:
  - "PostgreSQL предустановлен в sandbox-образе на localhost:5432 (пользователь postgres/postgres)"
  - "Каждый тест должен получать свою изолированную БД — создавай через CREATE DATABASE"
  - "httptest.NewServer запускает настоящий HTTP сервер для тестов"
  - "t.Cleanup регистрирует завершение сервера и удаление БД автоматически"
---

# Интеграционные тесты

Интеграционные тесты проверяют весь стек приложения как единое целое: реальная БД, реальный HTTP сервер, реальные запросы. Это защита от ошибок на стыке компонентов.

В отличие от unit-тестов с моками, интеграционные тесты:
- Находят ошибки несовместимости между слоями
- Тестируют SQL-запросы на реальной СУБД
- Проверяют сериализацию/десериализацию JSON end-to-end

В этом проекте PostgreSQL уже работает в sandbox-окружении на `localhost:5432` (пользователь и пароль — `postgres`). Твоя задача — для каждого теста поднять изолированную БД, собрать стек и вернуть HTTP-сервер.

## Что нужно сделать

Напиши helper `SetupTestServer`. Готовые сценарные тесты будут вызывать его и проверять стек через HTTP-запросы.

### Функция `SetupTestServer`

```go
package integration

func SetupTestServer(t *testing.T) *httptest.Server
```

Должна сделать:

1. **Подключиться к базовой БД и создать уникальную:**
   ```go
   const baseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

   base, _ := sql.Open("pgx", baseDSN)
   defer base.Close()

   dbName := fmt.Sprintf("t_%d_%d", time.Now().UnixNano(), rand.Intn(1<<30))
   base.Exec("CREATE DATABASE " + dbName)
   ```

2. **Подключиться к новой БД и мигрировать:**
   ```go
   testDSN := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
   database, _ := db.New(testDSN)
   db.Migrate(database)
   ```

3. **Собрать стек:**
   ```go
   repo := repository.New(database)
   h := handler.New(repo)
   router := server.NewRouter(h)
   fullHandler := middleware.Chain(router, middleware.Logger, middleware.Recovery)
   ```

4. **Запустить httptest и зарегистрировать cleanup:**
   ```go
   srv := httptest.NewServer(fullHandler)
   t.Cleanup(func() {
       srv.Close()
       database.Close()
       // удаляем тестовую БД
       cleanup, _ := sql.Open("pgx", baseDSN)
       defer cleanup.Close()
       cleanup.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '" + dbName + "'")
       cleanup.Exec("DROP DATABASE IF EXISTS " + dbName)
   })
   return srv
   ```

## Требования

- Функция **экспортируемая** (`SetupTestServer` с большой буквы) — вызывается из `scenarios_test.go`
- Каждый вызов создаёт **новую БД** с уникальным именем — тесты не должны видеть данные друг друга
- Cleanup обязательно: закрыть сервер, закрыть соединение, удалить БД
- Перед `DROP DATABASE` нужен `pg_terminate_backend` — иначе DROP не сработает из-за активных соединений
- Используй `t.Helper()` и `t.Fatalf` для понятного stack trace при ошибке

## Как это проверяется

Готовые сценарные тесты в том же пакете `integration` проверят твой helper через HTTP:

| Тест | Что проверяет |
|------|---------------|
| `TestHealthCheck` | GET /health → 200, body `{"status":"ok"}` |
| `TestCreateAndGetTask` | POST → 201, GET/{id} → 200 с теми же полями |
| `TestListTasks` | 3 создания, GET /tasks → массив из 3 |
| `TestUpdateTask` | POST → PUT → проверка изменений |
| `TestDeleteTask` | POST → DELETE → GET возвращает 404 |
| `TestValidation` | POST с пустым title → 400 |

Все шесть тестов вызывают твой `SetupTestServer(t)` для получения свежего сервера с чистой БД. Если хотя бы один тест видит данные другого — значит изоляция БД сломана.

## Зачем именно так, а не testcontainers

В реальных проектах часто используют [testcontainers-go](https://golang.testcontainers.org/) — он поднимает PostgreSQL в Docker для каждого теста. Здесь sandbox запускает тесты без Docker-доступа, поэтому используется уже работающая БД с изоляцией через отдельные базы. Идея та же — реальная PostgreSQL, изолированное окружение для каждого теста.
