package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"taskmanager/model"
)

// fakeRepo — in-memory реализация TaskRepository для тестов.
type fakeRepo struct {
	mu     sync.Mutex
	tasks  map[int64]model.Task
	nextID int64
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{tasks: make(map[int64]model.Task)}
}

func (f *fakeRepo) Create(_ context.Context, task model.Task) (model.Task, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextID++
	task.ID = f.nextID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	f.tasks[task.ID] = task
	return task, nil
}

func (f *fakeRepo) GetByID(_ context.Context, id int64) (model.Task, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	task, ok := f.tasks[id]
	if !ok {
		return model.Task{}, fmt.Errorf("task %d: not found", id)
	}
	return task, nil
}

func (f *fakeRepo) List(_ context.Context) ([]model.Task, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	tasks := make([]model.Task, 0, len(f.tasks))
	for _, t := range f.tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (f *fakeRepo) Update(_ context.Context, task model.Task) (model.Task, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.tasks[task.ID]; !ok {
		return model.Task{}, fmt.Errorf("task %d: not found", task.ID)
	}
	task.UpdatedAt = time.Now()
	f.tasks[task.ID] = task
	return task, nil
}

func (f *fakeRepo) Delete(_ context.Context, id int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.tasks[id]; !ok {
		return fmt.Errorf("task %d: not found", id)
	}
	delete(f.tasks, id)
	return nil
}

func TestCreateReturns201(t *testing.T) {
	h := New(newFakeRepo())
	body := strings.NewReader(`{"title":"Buy milk","description":"From the store"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var task model.Task
	if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if task.ID == 0 {
		t.Error("response task.ID = 0, want non-zero")
	}
	if task.Title != "Buy milk" {
		t.Errorf("response task.Title = %q, want %q", task.Title, "Buy milk")
	}
	if task.Description != "From the store" {
		t.Errorf("response task.Description = %q, want %q", task.Description, "From the store")
	}
	if task.Status != model.StatusPending {
		t.Errorf("response task.Status = %q, want %q (Create должен всегда ставить pending)",
			task.Status, model.StatusPending)
	}
}

func TestCreateIgnoresClientStatus(t *testing.T) {
	h := New(newFakeRepo())
	body := strings.NewReader(`{"title":"Test","status":"done"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()

	h.Create(w, req)

	var task model.Task
	if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if task.Status != model.StatusPending {
		t.Errorf("Status = %q, want %q (клиент не может задать начальный статус)",
			task.Status, model.StatusPending)
	}
}

func TestCreateEmptyTitleReturns400(t *testing.T) {
	h := New(newFakeRepo())
	body := strings.NewReader(`{"title":""}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateInvalidJSONReturns400(t *testing.T) {
	h := New(newFakeRepo())
	body := strings.NewReader(`not-json`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetByIDReturns200(t *testing.T) {
	repo := newFakeRepo()
	created, _ := repo.Create(context.Background(), model.NewTask("Test task", "Test desc"))

	h := New(repo)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", created.ID), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", created.ID))
	w := httptest.NewRecorder()

	h.GetByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got model.Task
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %d, want %d", got.ID, created.ID)
	}
	if got.Title != "Test task" {
		t.Errorf("Title = %q, want %q", got.Title, "Test task")
	}
	if got.Description != "Test desc" {
		t.Errorf("Description = %q, want %q", got.Description, "Test desc")
	}
}

func TestGetByIDInvalidIDReturns400(t *testing.T) {
	h := New(newFakeRepo())
	req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.GetByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d for non-numeric id", w.Code, http.StatusBadRequest)
	}
}

func TestGetByIDNotFoundReturns404(t *testing.T) {
	h := New(newFakeRepo())
	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	h.GetByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestListEmptyReturnsArray(t *testing.T) {
	h := New(newFakeRepo())
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if strings.TrimSpace(w.Body.String()) == "null" {
		t.Error("body = null, want []  (пустой список — это пустой массив)")
	}

	var tasks []model.Task
	if err := json.NewDecoder(w.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode: %v — body should be JSON array, not null", err)
	}
	if tasks == nil {
		t.Error("List returned null, want empty array []")
	}
	if len(tasks) != 0 {
		t.Errorf("len(tasks) = %d, want 0", len(tasks))
	}
}

func TestListReturnsAllTasks(t *testing.T) {
	repo := newFakeRepo()
	for _, title := range []string{"Task A", "Task B", "Task C"} {
		repo.Create(context.Background(), model.NewTask(title, ""))
	}
	h := New(repo)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var tasks []model.Task
	if err := json.NewDecoder(w.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("len(tasks) = %d, want 3", len(tasks))
	}
}

func TestDeleteReturns204(t *testing.T) {
	repo := newFakeRepo()
	task, _ := repo.Create(context.Background(), model.NewTask("To delete", ""))

	h := New(repo)
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", task.ID), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", task.ID))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
	if w.Body.Len() != 0 {
		t.Errorf("body = %q, want empty for 204", w.Body.String())
	}

	if _, err := repo.GetByID(context.Background(), task.ID); err == nil {
		t.Error("задача не удалена из хранилища после Delete")
	}
}

func TestDeleteNotFoundReturns404(t *testing.T) {
	h := New(newFakeRepo())
	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUpdateReturns200(t *testing.T) {
	repo := newFakeRepo()
	task, _ := repo.Create(context.Background(), model.NewTask("Old title", "Old desc"))

	h := New(repo)
	body := strings.NewReader(`{"title":"New title","description":"New desc","status":"in_progress"}`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", task.ID), body)
	req.SetPathValue("id", fmt.Sprintf("%d", task.ID))
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var updated model.Task
	if err := json.NewDecoder(w.Body).Decode(&updated); err != nil {
		t.Fatalf("decode: %v", err)
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

	stored, _ := repo.GetByID(context.Background(), task.ID)
	if stored.Title != "New title" {
		t.Errorf("в репозитории Title = %q, want %q (Update не применился к хранилищу)",
			stored.Title, "New title")
	}
}

func TestUpdateNotFoundReturns404(t *testing.T) {
	h := New(newFakeRepo())
	body := strings.NewReader(`{"title":"X"}`)
	req := httptest.NewRequest(http.MethodPut, "/tasks/999", body)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
