package server

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

	"taskmanager/handler"
	"taskmanager/model"
)

// fakeRepo — in-memory реализация для тестов роутера.
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

func newTestRouter(t *testing.T) (http.Handler, *fakeRepo) {
	t.Helper()
	repo := newFakeRepo()
	h := handler.New(repo)
	router := NewRouter(h)
	if router == nil {
		t.Fatal("NewRouter вернул nil")
	}
	return router, repo
}

// =====================  HEALTH  =====================

func TestHealthCheckStatus(t *testing.T) {
	router, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHealthCheckBody(t *testing.T) {
	router, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode /health: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf(`body["status"] = %q, want %q`, body["status"], "ok")
	}
}

func TestHealthCheckContentType(t *testing.T) {
	router, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if !strings.Contains(strings.ToLower(ct), "application/json") {
		t.Errorf("Content-Type = %q, want содержит %q", ct, "application/json")
	}
}

// =====================  ROUTES EXIST =====================

func TestPostTasksRoute(t *testing.T) {
	router, _ := newTestRouter(t)
	body := strings.NewReader(`{"title":"X"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("POST /tasks: status = %d, want %d (роут не зарегистрирован?)",
			w.Code, http.StatusCreated)
	}
}

func TestGetTasksRoute(t *testing.T) {
	router, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /tasks: status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetTasksByIDRoute(t *testing.T) {
	router, repo := newTestRouter(t)
	task, _ := repo.Create(context.Background(), model.NewTask("t", ""))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%d", task.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /tasks/{id}: status = %d, want %d (используется ли r.PathValue?)",
			w.Code, http.StatusOK)
	}
}

func TestPutTasksByIDRoute(t *testing.T) {
	router, repo := newTestRouter(t)
	task, _ := repo.Create(context.Background(), model.NewTask("old", ""))

	body := strings.NewReader(`{"title":"new"}`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", task.ID), body)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("PUT /tasks/{id}: status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestDeleteTasksByIDRoute(t *testing.T) {
	router, repo := newTestRouter(t)
	task, _ := repo.Create(context.Background(), model.NewTask("t", ""))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%d", task.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("DELETE /tasks/{id}: status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

// =====================  METHOD ROUTING =====================

func TestWrongMethodOnHealth(t *testing.T) {
	router, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /health: status = %d, want %d (Go 1.22 ServeMux должен возвращать 405 на несовпадающий метод)",
			w.Code, http.StatusMethodNotAllowed)
	}
}

func TestWrongMethodOnTasks(t *testing.T) {
	router, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodPatch, "/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("PATCH /tasks: status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

// =====================  404 =====================

func TestUnknownRouteReturns404(t *testing.T) {
	router, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/unknown/path", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestRootReturns404(t *testing.T) {
	router, _ := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Errorf("GET / вернул 200 — корневой роут не должен совпадать с зарегистрированными")
	}
}
