package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewTaskTitleAndDescription(t *testing.T) {
	task := NewTask("Купить продукты", "Молоко и хлеб")
	if task.Title != "Купить продукты" {
		t.Errorf("Title = %q, want %q", task.Title, "Купить продукты")
	}
	if task.Description != "Молоко и хлеб" {
		t.Errorf("Description = %q, want %q", task.Description, "Молоко и хлеб")
	}
}

func TestNewTaskDistinctTitleAndDescription(t *testing.T) {
	task := NewTask("title-value", "description-value")
	if task.Title == task.Description {
		t.Errorf("Title и Description оба = %q — поля перепутаны местами в конструкторе", task.Title)
	}
}

func TestNewTaskEmptyDescription(t *testing.T) {
	task := NewTask("only title", "")
	if task.Title != "only title" {
		t.Errorf("Title = %q, want %q", task.Title, "only title")
	}
	if task.Description != "" {
		t.Errorf("Description = %q, want empty string", task.Description)
	}
}

func TestNewTaskStatus(t *testing.T) {
	task := NewTask("Купить продукты", "Молоко и хлеб")
	if task.Status != StatusPending {
		t.Errorf("Status = %q, want %q", task.Status, StatusPending)
	}
}

func TestNewTaskZeroID(t *testing.T) {
	task := NewTask("Test", "Desc")
	if task.ID != 0 {
		t.Errorf("ID = %d, want 0 (новая задача без ID — его проставит база)", task.ID)
	}
}

func TestNewTaskCreatedAt(t *testing.T) {
	before := time.Now()
	task := NewTask("Test", "")
	after := time.Now()

	if task.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero, want non-zero time")
	}
	if task.CreatedAt.Before(before) || task.CreatedAt.After(after) {
		t.Errorf("CreatedAt = %v, want between %v and %v", task.CreatedAt, before, after)
	}
}

func TestNewTaskUpdatedAt(t *testing.T) {
	before := time.Now()
	task := NewTask("Test", "Desc")
	after := time.Now()

	if task.UpdatedAt.IsZero() {
		t.Error("UpdatedAt is zero, want non-zero time")
	}
	if task.UpdatedAt.Before(before) || task.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt = %v, want between %v and %v", task.UpdatedAt, before, after)
	}
}

func TestStatusConstants(t *testing.T) {
	cases := []struct {
		got, want Status
	}{
		{StatusPending, "pending"},
		{StatusInProgress, "in_progress"},
		{StatusDone, "done"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("Status = %q, want %q", c.got, c.want)
		}
	}
}

func TestValidateEmptyTitle(t *testing.T) {
	task := Task{Title: ""}
	err := task.Validate()
	if err == nil {
		t.Fatal("Validate() = nil, want error for empty title")
	}
	if err.Error() != "title is required" {
		t.Errorf("error message = %q, want %q", err.Error(), "title is required")
	}
}

func TestValidateNonEmptyTitle(t *testing.T) {
	task := Task{Title: "Buy milk"}
	if err := task.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestTaskJSONTags(t *testing.T) {
	task := Task{
		ID:          42,
		Title:       "Test task",
		Description: "Some description",
		Status:      StatusDone,
		CreatedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	fields := []string{"id", "title", "description", "status", "created_at", "updated_at"}
	for _, f := range fields {
		if _, ok := decoded[f]; !ok {
			t.Errorf("JSON missing field %q", f)
		}
	}
}

func TestTaskJSONRoundTrip(t *testing.T) {
	original := Task{
		ID:          1,
		Title:       "Write tests",
		Description: "Cover all handlers",
		Status:      StatusInProgress,
		CreatedAt:   time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 3, 2, 12, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	var decoded Task
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID = %d, want %d", decoded.ID, original.ID)
	}
	if decoded.Title != original.Title {
		t.Errorf("Title = %q, want %q", decoded.Title, original.Title)
	}
	if decoded.Description != original.Description {
		t.Errorf("Description = %q, want %q", decoded.Description, original.Description)
	}
	if decoded.Status != original.Status {
		t.Errorf("Status = %q, want %q", decoded.Status, original.Status)
	}
}
