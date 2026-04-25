package integration

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"taskmanager/db"
	"taskmanager/handler"
	"taskmanager/middleware"
	"taskmanager/repository"
	"taskmanager/server"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const baseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func SetupTestServer(t *testing.T) *httptest.Server {
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
		t.Fatalf("create database %s: %v", dbName, err)
	}
	base.Close()

	testDSN := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
	database, err := db.New(testDSN)
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}

	if err := db.Migrate(database); err != nil {
		t.Fatalf("db.Migrate: %v", err)
	}

	repo := repository.New(database)
	h := handler.New(repo)
	router := server.NewRouter(h)
	fullHandler := middleware.Chain(router, middleware.Logger, middleware.Recovery)

	srv := httptest.NewServer(fullHandler)

	t.Cleanup(func() {
		srv.Close()
		database.Close()

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

	return srv
}
