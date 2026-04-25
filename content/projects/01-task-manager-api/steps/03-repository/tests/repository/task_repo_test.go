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

func TestCreateAndGetByID(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	task := model.NewTask("Deploy app", "Set up CI/CD")
	task.Status = model.StatusPending

	created, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Error("created.ID = 0, want non-zero — RETURNING id не сработал")
	}
	if created.CreatedAt.IsZero() {
		t.Error("created.CreatedAt is zero — RETURNING created_at не сработал")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("created.UpdatedAt is zero — RETURNING updated_at не сработал")
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
}

func TestGetByIDNotFound(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 99999)
	if err == nil {
		t.Error("GetByID(nonexistent) = nil, want error")
	}
}

func TestListEmpty(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if tasks == nil {
		t.Error("List() = nil, want empty slice (not nil)")
	}
	if len(tasks) != 0 {
		t.Errorf("List() len = %d, want 0", len(tasks))
	}
}

func TestListReturnsAll(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	titles := []string{"Task A", "Task B", "Task C"}
	for _, title := range titles {
		if _, err := repo.Create(ctx, model.NewTask(title, "desc-"+title)); err != nil {
			t.Fatalf("Create %q: %v", title, err)
		}
		time.Sleep(2 * time.Millisecond)
	}

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("List() len = %d, want 3", len(tasks))
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

	if !tasks[0].CreatedAt.After(tasks[len(tasks)-1].CreatedAt) &&
		!tasks[0].CreatedAt.Equal(tasks[len(tasks)-1].CreatedAt) {
		t.Errorf("List должен сортировать ORDER BY created_at DESC: первая=%v, последняя=%v",
			tasks[0].CreatedAt, tasks[len(tasks)-1].CreatedAt)
	}
}

func TestUpdate(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	task := model.NewTask("Old title", "Old desc")
	created, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	originalUpdatedAt := created.UpdatedAt
	time.Sleep(2 * time.Millisecond)

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
	if !updated.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt не обновился: было %v, стало %v", originalUpdatedAt, updated.UpdatedAt)
	}

	stored, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID after Update: %v", err)
	}
	if stored.Title != "New title" {
		t.Errorf("в БД Title = %q, want %q (Update не сохранился)", stored.Title, "New title")
	}
	if stored.Status != model.StatusInProgress {
		t.Errorf("в БД Status = %q, want %q", stored.Status, model.StatusInProgress)
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

func TestDelete(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	task, err := repo.Create(ctx, model.NewTask("To delete", ""))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := repo.Delete(ctx, task.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = repo.GetByID(ctx, task.ID)
	if err == nil {
		t.Error("GetByID after Delete = nil, want error")
	}
}

func TestDeleteNotFound(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	if err := repo.Delete(ctx, 99999); err == nil {
		t.Error("Delete(nonexistent) = nil, want error")
	}
}
