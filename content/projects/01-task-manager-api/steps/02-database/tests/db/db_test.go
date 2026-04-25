package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const baseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

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

// ===================== NEW =====================

func TestNewConnects(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if database == nil {
		t.Fatal("New вернул nil *sql.DB без ошибки")
	}
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

func TestNewMalformedDSN(t *testing.T) {
	_, err := New("это вообще не URL")
	if err == nil {
		t.Error("New(malformed) = nil, want error")
	}
}

// ===================== MIGRATE =====================

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
		t.Error("table 'tasks' не существует после Migrate")
	}
}

func TestMigrateIdempotent(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	for i := 1; i <= 3; i++ {
		if err := Migrate(database); err != nil {
			t.Fatalf("Migrate #%d: %v (использовал ли ты CREATE TABLE IF NOT EXISTS?)", i, err)
		}
	}
}

func TestMigrateAllColumnsExist(t *testing.T) {
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
		"id": false, "title": false, "description": false,
		"status": false, "created_at": false, "updated_at": false,
	}
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if _, ok := want[col]; ok {
			want[col] = true
		}
	}
	for col, found := range want {
		if !found {
			t.Errorf("колонка %q отсутствует в таблице tasks", col)
		}
	}
}

func TestMigrateColumnTypes(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	want := map[string]string{
		"id":          "integer",
		"title":       "text",
		"description": "text",
		"status":      "text",
		"created_at":  "timestamp with time zone",
		"updated_at":  "timestamp with time zone",
	}

	rows, err := database.Query(
		`SELECT column_name, data_type FROM information_schema.columns
		 WHERE table_name = 'tasks'`,
	)
	if err != nil {
		t.Fatalf("query types: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, dtype string
		if err := rows.Scan(&name, &dtype); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if expected, ok := want[name]; ok && dtype != expected {
			t.Errorf("колонка %q: тип = %q, want %q", name, dtype, expected)
		}
	}
}

func TestMigrateTitleNotNull(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	_, err = database.Exec(`INSERT INTO tasks (title) VALUES (NULL)`)
	if err == nil {
		t.Error("INSERT с NULL title не упал — title должен быть NOT NULL")
		return
	}
	if !strings.Contains(strings.ToLower(err.Error()), "null") {
		t.Errorf("ошибка = %v, want содержащую 'null'", err)
	}
}

func TestMigrateDefaultStatus(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var status string
	err = database.QueryRow(
		`INSERT INTO tasks (title) VALUES ('test') RETURNING status`,
	).Scan(&status)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if status != "pending" {
		t.Errorf("status default = %q, want %q", status, "pending")
	}
}

func TestMigrateDefaultDescription(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var desc string
	err = database.QueryRow(
		`INSERT INTO tasks (title) VALUES ('t') RETURNING description`,
	).Scan(&desc)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if desc != "" {
		t.Errorf("description default = %q, want %q", desc, "")
	}
}

func TestMigrateDefaultTimestamps(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var createdAt, updatedAt time.Time
	err = database.QueryRow(
		`INSERT INTO tasks (title) VALUES ('t') RETURNING created_at, updated_at`,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	if createdAt.IsZero() {
		t.Error("created_at не получил значения по умолчанию (DEFAULT NOW())")
	}
	if updatedAt.IsZero() {
		t.Error("updated_at не получил значения по умолчанию (DEFAULT NOW())")
	}
}

func TestMigrateIDIsSerial(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var id1, id2 int64
	database.QueryRow(`INSERT INTO tasks (title) VALUES ('a') RETURNING id`).Scan(&id1)
	database.QueryRow(`INSERT INTO tasks (title) VALUES ('b') RETURNING id`).Scan(&id2)

	if id1 == 0 {
		t.Error("первый id = 0 — SERIAL не работает")
	}
	if id2 <= id1 {
		t.Errorf("id2 = %d, id1 = %d — id должны автоинкрементироваться", id2, id1)
	}
}

func TestMigrateIDIsPrimaryKey(t *testing.T) {
	dsn := uniqueDB(t)
	database, err := New(dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer database.Close()

	if err := Migrate(database); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var col string
	err = database.QueryRow(
		`SELECT a.attname
		 FROM pg_index i
		 JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		 WHERE i.indrelid = 'tasks'::regclass AND i.indisprimary`,
	).Scan(&col)
	if err != nil {
		t.Fatalf("query primary key: %v", err)
	}
	if col != "id" {
		t.Errorf("primary key = %q, want %q", col, "id")
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
		`INSERT INTO tasks (title, description, status) VALUES ('t', 'd', 'in_progress') RETURNING id`,
	).Scan(&id)
	if err != nil {
		t.Errorf("insert with all fields: %v", err)
	}
	if id == 0 {
		t.Error("inserted id = 0, want non-zero SERIAL value")
	}
}
