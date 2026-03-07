package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTodoFields(t *testing.T) {
	now := time.Now()
	todo := Todo{
		ID:          1,
		Title:       "Test",
		Description: "Description",
		Done:        true,
		CreatedAt:   now,
	}

	if todo.ID != 1 {
		t.Errorf("ID = %d, want 1", todo.ID)
	}
	if todo.Title != "Test" {
		t.Errorf("Title = %q, want %q", todo.Title, "Test")
	}
	if todo.Description != "Description" {
		t.Errorf("Description = %q, want %q", todo.Description, "Description")
	}
	if !todo.Done {
		t.Error("Done = false, want true")
	}
	if !todo.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", todo.CreatedAt, now)
	}
}

func TestTodoJSON(t *testing.T) {
	todo := Todo{
		ID:          1,
		Title:       "Buy milk",
		Description: "Go to the store",
		Done:        false,
		CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	fields := []string{"id", "title", "description", "done", "created_at"}
	for _, f := range fields {
		if _, ok := decoded[f]; !ok {
			t.Errorf("JSON missing field %q", f)
		}
	}
}

func TestTodoJSONRoundTrip(t *testing.T) {
	original := Todo{
		ID:          42,
		Title:       "Test task",
		Description: "Some description",
		Done:        true,
		CreatedAt:   time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var decoded Todo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %d, want %d", decoded.ID, original.ID)
	}
	if decoded.Title != original.Title {
		t.Errorf("Title = %q, want %q", decoded.Title, original.Title)
	}
	if decoded.Done != original.Done {
		t.Errorf("Done = %v, want %v", decoded.Done, original.Done)
	}
}

func TestCreateTodoRequest(t *testing.T) {
	body := `{"title":"New task","description":"Task description"}`

	var req CreateTodoRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if req.Title != "New task" {
		t.Errorf("Title = %q, want %q", req.Title, "New task")
	}
	if req.Description != "Task description" {
		t.Errorf("Description = %q, want %q", req.Description, "Task description")
	}
}

func TestUpdateTodoRequest(t *testing.T) {
	body := `{"title":"Updated","done":true}`

	var req UpdateTodoRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if req.Title == nil || *req.Title != "Updated" {
		t.Errorf("Title = %v, want %q", req.Title, "Updated")
	}
	if req.Description != nil {
		t.Errorf("Description = %v, want nil", req.Description)
	}
	if req.Done == nil || *req.Done != true {
		t.Errorf("Done = %v, want true", req.Done)
	}
}

func TestUpdateTodoRequestAllNil(t *testing.T) {
	body := `{}`

	var req UpdateTodoRequest
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if req.Title != nil {
		t.Errorf("Title = %v, want nil", req.Title)
	}
	if req.Description != nil {
		t.Errorf("Description = %v, want nil", req.Description)
	}
	if req.Done != nil {
		t.Errorf("Done = %v, want nil", req.Done)
	}
}
