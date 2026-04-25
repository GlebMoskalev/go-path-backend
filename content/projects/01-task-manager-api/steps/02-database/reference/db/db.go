package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(connString string) (*sql.DB, error) {
	database, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return database, nil
}

func Migrate(database *sql.DB) error {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id          SERIAL PRIMARY KEY,
			title       TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			status      TEXT NOT NULL DEFAULT 'pending',
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}
