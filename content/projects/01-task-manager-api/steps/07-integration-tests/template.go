package integration

// Реализуйте SetupTestServer — helper для интеграционных тестов.
//
// PostgreSQL уже работает в sandbox-образе на localhost:5432 (postgres/postgres).
// Твоя задача — создать изолированную БД для теста и собрать полный стек.
//
// Импорты:
//   import (
//       "database/sql"
//       "fmt"
//       "math/rand"
//       "net/http/httptest"
//       "testing"
//       "time"
//
//       "taskmanager/db"
//       "taskmanager/handler"
//       "taskmanager/middleware"
//       "taskmanager/repository"
//       "taskmanager/server"
//
//       _ "github.com/jackc/pgx/v5/stdlib"
//   )
//
// const baseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
//
// func SetupTestServer(t *testing.T) *httptest.Server
//
// Шаги:
//   1. t.Helper()
//   2. base, _ := sql.Open("pgx", baseDSN); defer base.Close()
//   3. Создай уникальную БД:
//      dbName := fmt.Sprintf("t_%d_%d", time.Now().UnixNano(), rand.Intn(1<<30))
//      base.Exec("CREATE DATABASE " + dbName)
//   4. testDSN := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
//      database, _ := db.New(testDSN)
//      db.Migrate(database)
//   5. Собери цепочку: repository → handler → router → middleware.Chain
//   6. srv := httptest.NewServer(fullHandler)
//   7. t.Cleanup:
//      - srv.Close()
//      - database.Close()
//      - открой baseDSN и сделай pg_terminate_backend, потом DROP DATABASE
//   8. return srv
//
// ВАЖНО: функция экспортируемая (SetupTestServer с большой буквы) —
// её вызывают из scenarios_test.go в этом же пакете.
