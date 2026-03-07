package storage

import (
	"sort"
	"testing"

	"todoapi/model"
)

func TestCreate(t *testing.T) {
	s := NewTodoStorage()

	todo := s.Create(model.Todo{Title: "Test", Description: "Desc"})

	if todo.ID != 1 {
		t.Errorf("ID = %d, want 1", todo.ID)
	}
	if todo.Title != "Test" {
		t.Errorf("Title = %q, want %q", todo.Title, "Test")
	}
	if todo.Description != "Desc" {
		t.Errorf("Description = %q, want %q", todo.Description, "Desc")
	}
	if todo.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestCreateAutoIncrement(t *testing.T) {
	s := NewTodoStorage()

	t1 := s.Create(model.Todo{Title: "First"})
	t2 := s.Create(model.Todo{Title: "Second"})

	if t1.ID != 1 {
		t.Errorf("first ID = %d, want 1", t1.ID)
	}
	if t2.ID != 2 {
		t.Errorf("second ID = %d, want 2", t2.ID)
	}
}

func TestGetByID(t *testing.T) {
	s := NewTodoStorage()
	created := s.Create(model.Todo{Title: "Test"})

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Title != "Test" {
		t.Errorf("Title = %q, want %q", got.Title, "Test")
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewTodoStorage()

	_, err := s.GetByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestGetAll(t *testing.T) {
	s := NewTodoStorage()
	s.Create(model.Todo{Title: "First"})
	s.Create(model.Todo{Title: "Second"})

	todos := s.GetAll()
	if len(todos) != 2 {
		t.Fatalf("len = %d, want 2", len(todos))
	}

	sort.Slice(todos, func(i, j int) bool { return todos[i].ID < todos[j].ID })
	if todos[0].Title != "First" {
		t.Errorf("todos[0].Title = %q, want %q", todos[0].Title, "First")
	}
	if todos[1].Title != "Second" {
		t.Errorf("todos[1].Title = %q, want %q", todos[1].Title, "Second")
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewTodoStorage()

	todos := s.GetAll()
	if todos == nil {
		t.Fatal("GetAll returned nil, want empty slice")
	}
	if len(todos) != 0 {
		t.Errorf("len = %d, want 0", len(todos))
	}
}

func TestUpdate(t *testing.T) {
	s := NewTodoStorage()
	created := s.Create(model.Todo{Title: "Original"})

	created.Title = "Updated"
	created.Done = true

	updated, err := s.Update(created.ID, created)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Title != "Updated" {
		t.Errorf("Title = %q, want %q", updated.Title, "Updated")
	}
	if !updated.Done {
		t.Error("Done = false, want true")
	}
}

func TestUpdateNotFound(t *testing.T) {
	s := NewTodoStorage()

	_, err := s.Update(999, model.Todo{})
	if err != ErrNotFound {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestDelete(t *testing.T) {
	s := NewTodoStorage()
	created := s.Create(model.Todo{Title: "Test"})

	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Errorf("after delete: error = %v, want ErrNotFound", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	s := NewTodoStorage()

	err := s.Delete(999)
	if err != ErrNotFound {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}
