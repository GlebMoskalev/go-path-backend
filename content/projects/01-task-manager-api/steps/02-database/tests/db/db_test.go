package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const baseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

// uniqueDB создаёт изолированную БД для теста.
// PostgreSQL должен быть запущен на localhost:5432 (предустановлен в sandbox-образе).
// Возвращает DSN новой БД и регистрирует cleanup для её удаления.
func uniqueDB(t *testing.T) string {
	t.Helper()

	base, err := sql.Open("pgx", baseDSN)
	if err != nil {
		t.Fatalf("open base db: %v", err)
	}
	if err := base.Ping(); err != nil {
		base.Close()
		t.Fatalf("PostgreSQL недоступен на localhost:5432: %v", err)
	}

	dbName := fmt.Sprintf("t_%d_%d", time.Now().UnixNano(), rand.Intn(1<<30))
	if _, err := base.Exec("CREATE DATABASE " + dbName); err != nil {
		base.Close()
		t.Fatalf("create db %s: %v", dbName, err)
	}
	base.Close()

	t.Cleanup(func() {
		cleanup, err := sql.Open("pgx", baseDSN)
		if err != nil {
			return
		}
		defer cleanup.Close()
		cleanup.Exec(fmt.Sprintf(
			"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s'",
			dbName))
		cleanup.Exec("DROP DATABASE IF EXISTS " + dbName)
	})

	return fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
}

func TestNewConnects(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		t.Errorf("Ping after New: %v", err)
	}
}

func TestNewInvalidConnString(t *testing.T) {
	_, err := New("postgres://nouser:nopass@127.0.0.1:9/nodb?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Error("New(invalid) = nil, want error")
	}
}

func TestMigrateCreatesTable(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var exists bool
	err = database.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_name = 'tasks'
		)`,
	).Scan(&exists)
	if err != nil {
		t.Fatalf("check table exists: %v", err)
	}
	if !exists {
		t.Error("table 'tasks' does not exist after Migrate")
	}
}

func TestMigrateIdempotent(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("first Migrate: %v", err)
	}
	if err := Migrate(database); err != nil {
		t.Errorf("second Migrate: %v, want nil (idempotent)", err)
	}
}

func TestMigrateTableColumns(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	rows, err := database.Query(
		`SELECT column_name FROM information_schema.columns
		 WHERE table_name = 'tasks' ORDER BY ordinal_position`,
	)
	if err != nil {
		t.Fatalf("query columns: %v", err)
	}
	defer rows.Close()

	want := map[string]bool{
		"id": true, "title": true, "description": true,
		"status": true, "created_at": true, "updated_at": true,
	}
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			t.Fatalf("scan: %v", err)
		}
		delete(want, col)
	}
	for col := range want {
		t.Errorf("column %q missing in tasks table", col)
	}
}

func TestInsertAfterMigrate(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var id int64
	err = database.QueryRow(
		`INSERT INTO tasks (title) VALUES ('test') RETURNING id`,
	).Scan(&id)
	if err != nil {
		t.Errorf("insert after migrate: %v", err)
	}
	if id == 0 {
		t.Error("inserted id = 0, want non-zero SERIAL value")
	}
}
