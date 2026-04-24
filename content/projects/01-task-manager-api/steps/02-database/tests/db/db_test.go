package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgres(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}
	return connStr
}

func TestNewConnects(t *testing.T) {
	connStr := setupPostgres(t)
	db, err := New(connStr)
	if err != nil {
		t.Fatalf("New(%q) = %v, want nil error", connStr, err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Errorf("Ping after New: %v", err)
	}
}

func TestNewInvalidConnString(t *testing.T) {
	_, err := New("postgres://invalid:invalid@localhost:9999/nodb?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Error("New(invalid) = nil, want error")
	}
}

func TestMigrateCreatesTable(t *testing.T) {
	connStr := setupPostgres(t)
	db, err := New(connStr)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var exists bool
	err = db.QueryRow(
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
	connStr := setupPostgres(t)
	db, err := New(connStr)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("first Migrate: %v", err)
	}
	if err := Migrate(db); err != nil {
		t.Errorf("second Migrate: %v, want nil (idempotent)", err)
	}
}

func TestMigrateTableColumns(t *testing.T) {
	connStr := setupPostgres(t)
	database, err := New(connStr)
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
	connStr := setupPostgres(t)
	database, err := New(connStr)
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

// подавляем "imported and not used" для sql если не используется напрямую
var _ = sql.ErrNoRows
