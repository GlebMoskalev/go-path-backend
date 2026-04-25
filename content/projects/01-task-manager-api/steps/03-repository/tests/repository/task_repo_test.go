package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"taskmanager/db"
	"taskmanager/model"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const baseDSN = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func setupTestRepo(t *testing.T) *TaskRepository {
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
		t.Fatalf("create db: %v", err)
	}
	base.Close()

	dsn := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)
	database, err := db.New(dsn)
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("db.Migrate: %v", err)
	}

	t.Cleanup(func() {
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

	return New(database)
}

// =====================  CREATE  =====================

func TestCreateAssignsID(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, model.NewTask("t", ""))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Error("created.ID = 0 — RETURNING id не сработал")
	}
}

func TestCreateAssignsTimestamps(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	created, err := repo.Create(ctx, model.NewTask("t", ""))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.CreatedAt.IsZero() {
		t.Error("created.CreatedAt is zero — RETURNING created_at не сработал")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("created.UpdatedAt is zero — RETURNING updated_at не сработал")
	}
}

func TestCreateUniqueIDs(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	a, _ := repo.Create(ctx, model.NewTask("a", ""))
	b, _ := repo.Create(ctx, model.NewTask("b", ""))

	if a.ID == b.ID {
		t.Errorf("оба ID = %d, want разные (SERIAL должен инкрементироваться)", a.ID)
	}
}

func TestCreateAndGetByID(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	task := model.NewTask("Deploy app", "Set up CI/CD")
	task.Status = model.StatusPending

	created, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID(%d): %v", created.ID, err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %d, want %d", got.ID, created.ID)
	}
	if got.Title != "Deploy app" {
		t.Errorf("Title = %q, want %q", got.Title, "Deploy app")
	}
	if got.Description != "Set up CI/CD" {
		t.Errorf("Description = %q, want %q", got.Description, "Set up CI/CD")
	}
	if got.Status != model.StatusPending {
		t.Errorf("Status = %q, want %q", got.Status, model.StatusPending)
	}
	if got.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero после GetByID")
	}
}

// =====================  GET BY ID  =====================

func TestGetByIDNotFound(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 99999)
	if err == nil {
		t.Error("GetByID(несуществующий) = nil, want error")
	}
}

func TestGetByIDNotFoundIsNotErrNoRows(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 99999)
	if err == sql.ErrNoRows {
		t.Error("GetByID не должен возвращать sql.ErrNoRows напрямую — заверни в fmt.Errorf")
	}
}

// =====================  LIST  =====================

func TestListEmptyReturnsEmptySlice(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if tasks == nil {
		t.Error("List() = nil, want []model.Task{} (НЕ nil — JSON-API клиенты ждут массив)")
	}
	if len(tasks) != 0 {
		t.Errorf("len = %d, want 0", len(tasks))
	}
}

func TestListReturnsAllTitles(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	titles := []string{"Task A", "Task B", "Task C"}
	for _, title := range titles {
		if _, err := repo.Create(ctx, model.NewTask(title, "desc-"+title)); err != nil {
			t.Fatalf("Create %q: %v", title, err)
		}
	}

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("len = %d, want 3", len(tasks))
	}

	got := map[string]bool{}
	for _, task := range tasks {
		got[task.Title] = true
	}
	for _, want := range titles {
		if !got[want] {
			t.Errorf("title %q отсутствует в результате List", want)
		}
	}
}

func TestListAllFieldsPopulated(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	repo.Create(ctx, model.NewTask("with-desc", "the description"))

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("len = %d, want 1", len(tasks))
	}
	got := tasks[0]
	if got.ID == 0 {
		t.Error("ID = 0 в результате List")
	}
	if got.Title != "with-desc" {
		t.Errorf("Title = %q", got.Title)
	}
	if got.Description != "the description" {
		t.Errorf("Description = %q", got.Description)
	}
	if got.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
}

func TestListOrderByCreatedAtDesc(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	for _, title := range []string{"first", "second", "third"} {
		repo.Create(ctx, model.NewTask(title, ""))
		time.Sleep(2 * time.Millisecond)
	}

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("len = %d, want 3", len(tasks))
	}

	if tasks[0].Title != "third" {
		t.Errorf("первая задача = %q, want %q (ORDER BY created_at DESC — новые сначала)",
			tasks[0].Title, "third")
	}
	if tasks[2].Title != "first" {
		t.Errorf("последняя задача = %q, want %q", tasks[2].Title, "first")
	}
}

// =====================  UPDATE  =====================

func TestUpdateAllFields(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	created, _ := repo.Create(ctx, model.NewTask("Old title", "Old desc"))

	created.Title = "New title"
	created.Description = "New desc"
	created.Status = model.StatusInProgress
	updated, err := repo.Update(ctx, created)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Title != "New title" {
		t.Errorf("Title = %q, want %q", updated.Title, "New title")
	}
	if updated.Description != "New desc" {
		t.Errorf("Description = %q, want %q", updated.Description, "New desc")
	}
	if updated.Status != model.StatusInProgress {
		t.Errorf("Status = %q, want %q", updated.Status, model.StatusInProgress)
	}
}

func TestUpdatePersists(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	created, _ := repo.Create(ctx, model.NewTask("Old", ""))

	created.Title = "New"
	created.Status = model.StatusDone
	if _, err := repo.Update(ctx, created); err != nil {
		t.Fatalf("Update: %v", err)
	}

	stored, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if stored.Title != "New" {
		t.Errorf("в БД Title = %q, want %q (Update не сохранился)", stored.Title, "New")
	}
	if stored.Status != model.StatusDone {
		t.Errorf("в БД Status = %q, want %q", stored.Status, model.StatusDone)
	}
}

func TestUpdateChangesUpdatedAt(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	created, _ := repo.Create(ctx, model.NewTask("t", ""))
	originalUpdatedAt := created.UpdatedAt
	time.Sleep(5 * time.Millisecond)

	created.Title = "new"
	updated, err := repo.Update(ctx, created)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !updated.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt не обновился: было %v, стало %v (нужно SET updated_at=NOW())",
			originalUpdatedAt, updated.UpdatedAt)
	}
}

func TestUpdateNotFound(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	ghost := model.Task{ID: 99999, Title: "ghost", Status: model.StatusPending}
	if _, err := repo.Update(ctx, ghost); err == nil {
		t.Error("Update(несуществующий) = nil, want error")
	}
}

// =====================  DELETE  =====================

func TestDelete(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	task, _ := repo.Create(ctx, model.NewTask("To delete", ""))

	if err := repo.Delete(ctx, task.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := repo.GetByID(ctx, task.ID)
	if err == nil {
		t.Error("GetByID после Delete = nil, want error (задача должна быть удалена)")
	}
}

func TestDeleteNotFound(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	if err := repo.Delete(ctx, 99999); err == nil {
		t.Error("Delete(несуществующий) = nil, want error (нужно проверять RowsAffected)")
	}
}

func TestDeleteOnlyTargetTask(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	a, _ := repo.Create(ctx, model.NewTask("Keep", ""))
	b, _ := repo.Create(ctx, model.NewTask("Delete", ""))

	if err := repo.Delete(ctx, b.ID); err != nil {
		t.Fatalf("Delete(b): %v", err)
	}

	if _, err := repo.GetByID(ctx, a.ID); err != nil {
		t.Errorf("Delete удалил не ту задачу: %v", err)
	}
}

// =====================  CONTEXT =====================

func TestCreateRespectsCancelledContext(t *testing.T) {
	repo := setupTestRepo(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := repo.Create(ctx, model.NewTask("x", "")); err == nil {
		t.Error("Create с отменённым контекстом = nil, want context.Canceled (передаёшь ли ты ctx в QueryRowContext?)")
	}
}
