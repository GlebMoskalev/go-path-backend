package db

// Реализуйте функции для работы с PostgreSQL.
//
// Импорты:
//   import (
//       "database/sql"
//       _ "github.com/jackc/pgx/v5/stdlib"
//   )
//
// 1. func New(connString string) (*sql.DB, error)
//    - sql.Open("pgx", connString) — открывает соединение
//    - db.Ping() — проверяет реальное подключение
//    - если Ping упал: db.Close() и верни ошибку
//
// 2. func Migrate(db *sql.DB) error
//    - выполни CREATE TABLE IF NOT EXISTS tasks (
//          id          SERIAL PRIMARY KEY,
//          title       TEXT NOT NULL,
//          description TEXT NOT NULL DEFAULT '',
//          status      TEXT NOT NULL DEFAULT 'pending',
//          created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
//          updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
//      )
