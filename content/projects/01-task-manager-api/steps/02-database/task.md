---
title: "Подключение к PostgreSQL"
difficulty: easy
order: 2
file: "db/db.go"
hints:
  - "Используй pgx/v5/stdlib для совместимости с database/sql интерфейсом"
  - "Connection string формат: postgres://user:password@host:port/dbname?sslmode=disable"
  - "Всегда вызывай db.Ping() после открытия — убедись что соединение живое"
  - "В Migrate используй SERIAL для автоинкремента, а не AUTOINCREMENT как в SQLite"
  - "Тесты используют PostgreSQL на localhost:5432 (предустановлен в sandbox)"
---

# Подключение к PostgreSQL

PostgreSQL — самая популярная реляционная СУБД в Go-проектах. В этом шаге ты подключишься к базе данных и создашь схему таблицы.

Go работает с базами данных через универсальный интерфейс `database/sql`. Это значит, что код работы с запросами не зависит от конкретной СУБД — можно заменить PostgreSQL на MySQL и изменить лишь строку подключения и драйвер.

## Что нужно сделать

Реализуй два функции в пакете `db`.

### 1. Функция `New`

```go
func New(connString string) (*sql.DB, error)
```

- Открывает соединение через `sql.Open("pgx", connString)` — здесь `"pgx"` это имя драйвера
- Вызывает `db.Ping()` для проверки реального соединения (Open не устанавливает соединение!)
- Если Ping упал — закрой соединение и верни ошибку

### 2. Функция `Migrate`

```go
func Migrate(db *sql.DB) error
```

Выполняет миграцию — создаёт таблицу `tasks` если её ещё нет:

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
```

## Требования

- Импортируй драйвер через `_ "github.com/jackc/pgx/v5/stdlib"` (side-effect import для регистрации драйвера)
- `New` должна закрывать `*sql.DB` при ошибке Ping — не оставляй открытые соединения
- `Migrate` должна быть идемпотентной: повторный вызов не должен возвращать ошибку
- Используй `TIMESTAMPTZ` (timestamp with time zone) вместо `TIMESTAMP` — это стандарт для production

## Пример использования

```go
db, err := db.New("postgres://user:pass@localhost:5432/taskmanager?sslmode=disable")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

if err := db.Migrate(db); err != nil {
    log.Fatal(err)
}
// таблица tasks создана, готово к работе
```
