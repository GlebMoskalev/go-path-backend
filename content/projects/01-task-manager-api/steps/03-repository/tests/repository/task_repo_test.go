package repository

import (
	"context"
	"testing"
	"time"

	"taskmanager/db"
	"taskmanager/model"

	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *TaskRepository {
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
		t.Fatalf("start postgres: %v", err)
	}
	t.Cleanup(func() { container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	database, err := db.New(connStr)
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	if err := db.Migrate(database); err != nil {
		t.Fatalf("db.Migrate: %v", err)
	}

	return New(database)
}

func TestCreateAndGetByID(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	task := model.NewTask("Deploy app", "Set up CI/CD")
	task.Status = model.StatusPending

	created, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Error("created.ID = 0, want non-zero")
	}
	if created.CreatedAt.IsZero() {
		t.Error("created.CreatedAt is zero")
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID(%d): %v", created.ID, err)
	}
	if got.Title != task.Title {
		t.Errorf("Title = %q, want %q", got.Title, task.Title)
	}
	if got.Status != model.StatusPending {
		t.Errorf("Status = %q, want %q", got.Status, model.StatusPending)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 99999)
	if err == nil {
		t.Error("GetByID(nonexistent) = nil, want error")
	}
}

func TestListEmpty(t *testing.T) {
	repo := setupTestDB(t)
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
	repo := setupTestDB(t)
	ctx := context.Background()

	for _, title := range []string{"Task A", "Task B", "Task C"} {
		if _, err := repo.Create(ctx, model.NewTask(title, "")); err != nil {
			t.Fatalf("Create %q: %v", title, err)
		}
	}

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("List() len = %d, want 3", len(tasks))
	}
}

func TestUpdate(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	task := model.NewTask("Old title", "Old desc")
	created, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	created.Title = "New title"
	created.Status = model.StatusInProgress
	updated, err := repo.Update(ctx, created)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Title != "New title" {
		t.Errorf("Title = %q, want %q", updated.Title, "New title")
	}
	if updated.Status != model.StatusInProgress {
		t.Errorf("Status = %q, want %q", updated.Status, model.StatusInProgress)
	}
	if updated.UpdatedAt.IsZero() {
		t.Error("UpdatedAt is zero after update")
	}
}

func TestDelete(t *testing.T) {
	repo := setupTestDB(t)
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
	repo := setupTestDB(t)
	ctx := context.Background()

	if err := repo.Delete(ctx, 99999); err == nil {
		t.Error("Delete(nonexistent) = nil, want error")
	}
}
