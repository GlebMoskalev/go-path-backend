package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func newTestRouter() http.Handler {
	h := handler.New(newFakeRepo())
	return NewRouter(h)
}

func TestHealthCheck(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode /health response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want %q", body["status"], "ok")
	}
}

func TestHealthCheckContentType(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if ct == "" {
		t.Error("Content-Type header missing for /health")
	}
}

func TestUnknownRouteReturns404(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/unknown/path", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestTasksRouteExists(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Error("GET /tasks returned 404, route should be registered")
	}
}
