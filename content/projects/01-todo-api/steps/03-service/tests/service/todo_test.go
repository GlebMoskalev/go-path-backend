package service

import (
	"testing"

	"todoapi/model"
	"todoapi/storage"
)

func newTestService() *TodoService {
	return NewTodoService(storage.NewTodoStorage())
}

func TestServiceCreate(t *testing.T) {
	svc := newTestService()

	todo := svc.Create(model.CreateTodoRequest{
		Title:       "Test",
		Description: "Desc",
	})

	if todo.ID == 0 {
		t.Error("ID should not be 0")
	}
	if todo.Title != "Test" {
		t.Errorf("Title = %q, want %q", todo.Title, "Test")
	}
	if todo.Description != "Desc" {
		t.Errorf("Description = %q, want %q", todo.Description, "Desc")
	}
	if todo.Done {
		t.Error("Done should be false for new todo")
	}
	if todo.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestServiceGetByID(t *testing.T) {
	svc := newTestService()

	created := svc.Create(model.CreateTodoRequest{Title: "Test"})

	got, err := svc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Title != "Test" {
		t.Errorf("Title = %q, want %q", got.Title, "Test")
	}
}

func TestServiceGetByIDNotFound(t *testing.T) {
	svc := newTestService()

	_, err := svc.GetByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestServiceGetAll(t *testing.T) {
	svc := newTestService()

	svc.Create(model.CreateTodoRequest{Title: "First"})
	svc.Create(model.CreateTodoRequest{Title: "Second"})

	todos := svc.GetAll()
	if len(todos) != 2 {
		t.Errorf("len = %d, want 2", len(todos))
	}
}

func TestServiceUpdate(t *testing.T) {
	svc := newTestService()

	created := svc.Create(model.CreateTodoRequest{
		Title:       "Original",
		Description: "Desc",
	})

	newTitle := "Updated"
	done := true

	updated, err := svc.Update(created.ID, model.UpdateTodoRequest{
		Title: &newTitle,
		Done:  &done,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if updated.Title != "Updated" {
		t.Errorf("Title = %q, want %q", updated.Title, "Updated")
	}
	if updated.Description != "Desc" {
		t.Errorf("Description changed to %q, should stay %q", updated.Description, "Desc")
	}
	if !updated.Done {
		t.Error("Done = false, want true")
	}
}

func TestServiceUpdatePartial(t *testing.T) {
	svc := newTestService()

	created := svc.Create(model.CreateTodoRequest{
		Title:       "Original",
		Description: "Original desc",
	})

	newDesc := "New desc"
	updated, err := svc.Update(created.ID, model.UpdateTodoRequest{
		Description: &newDesc,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if updated.Title != "Original" {
		t.Errorf("Title = %q, want %q", updated.Title, "Original")
	}
	if updated.Description != "New desc" {
		t.Errorf("Description = %q, want %q", updated.Description, "New desc")
	}
}

func TestServiceUpdateNotFound(t *testing.T) {
	svc := newTestService()

	title := "New"
	_, err := svc.Update(999, model.UpdateTodoRequest{Title: &title})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestServiceDelete(t *testing.T) {
	svc := newTestService()

	created := svc.Create(model.CreateTodoRequest{Title: "Test"})

	err := svc.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = svc.GetByID(created.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestServiceDeleteNotFound(t *testing.T) {
	svc := newTestService()

	err := svc.Delete(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
