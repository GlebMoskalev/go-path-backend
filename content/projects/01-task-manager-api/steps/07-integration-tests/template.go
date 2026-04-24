package integration

// Напишите интеграционные тесты для всего стека приложения.
//
// Импорты:
//   import (
//       "context"
//       "encoding/json"
//       "fmt"
//       "net/http"
//       "net/http/httptest"
//       "strings"
//       "testing"
//       "time"
//
//       "taskmanager/db"
//       "taskmanager/handler"
//       "taskmanager/middleware"
//       "taskmanager/model"
//       "taskmanager/repository"
//       "taskmanager/server"
//
//       testcontainers "github.com/testcontainers/testcontainers-go"
//       tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
//       "github.com/testcontainers/testcontainers-go/wait"
//   )
//
// 1. func setupTestServer(t *testing.T) *httptest.Server
//    - запускает PostgreSQL через testcontainers
//    - создаёт db, запускает Migrate
//    - собирает repository → handler → router → middleware.Chain
//    - возвращает httptest.NewServer
//    - регистрирует t.Cleanup для остановки контейнера и сервера
//
// 2. Реализуй 6 тестов:
//    - TestHealthCheck         — GET /health → 200
//    - TestCreateAndGetTask    — POST создаёт, GET/{id} возвращает
//    - TestListTasks           — создать несколько, GET /tasks возвращает все
//    - TestUpdateTask          — создать, PUT обновить, проверить
//    - TestDeleteTask          — создать, DELETE, GET → 404
//    - TestValidation          — POST с пустым title → 400
