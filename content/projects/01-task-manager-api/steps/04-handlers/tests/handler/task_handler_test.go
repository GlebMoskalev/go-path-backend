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
	failOn string // если != "" — методы возвращают ошибку
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{tasks: make(map[int64]model.Task)}
}

func (f *fakeRepo) Create(_ context.Context, task model.Task) (model.Task, error) {
	if f.failOn == "Create" {
		return model.Task{}, fmt.Errorf("simulated failure")
	}
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
	if f.failOn == "List" {
		return nil, fmt.Errorf("simulated failure")
	}
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

// assertJSONContentType проверяет что ответ содержит Content-Type: application/json.
func assertJSONContentType(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	ct := w.Header().Get("Content-Type")
	if ct == "" {
		t.Errorf("Content-Type header отсутствует — клиент не поймёт что это JSON")
		return
	}
	if !strings.Contains(strings.ToLower(ct), "application/json") {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}
}

// assertErrorBody проверяет что тело ответа имеет вид {"error": "..."}.
func assertErrorBody(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("error response не валидный JSON: %v (body: %q)", err, w.Body.String())
	}
	msg, ok := body["error"]
	if !ok {
		t.Errorf("error response = %v, want объект с полем %q", body, "error")
		return
	}
	if msg == "" {
		t.Errorf("error response %q пустое — клиент не узнает что пошло не так", "error")
	}
}

// =====================  CREATE  =====================

func TestCreateReturns201(t *testing.T) {
	h := New(newFakeRepo())
	body := strings.NewReader(`{"title":"Buy milk","description":"From the store"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	assertJSONContentType(t, w)

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

func TestCreatePersistsToRepo(t *testing.T) {
	repo := newFakeRepo()
	h := New(repo)

	body := strings.NewReader(`{"title":"Persisted","description":"Check repo"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()
	h.Create(w, req)

	var resp model.Task
	json.NewDecoder(w.Body).Decode(&resp)

	stored, err := repo.GetByID(context.Background(), resp.ID)
	if err != nil {
		t.Fatalf("задача не сохранилась в репозитории: %v", err)
	}
	if stored.Title != "Persisted" {
		t.Errorf("в репо Title = %q, want %q", stored.Title, "Persisted")
	}
	if stored.Description != "Check repo" {
		t.Errorf("в репо Description = %q, want %q", stored.Description, "Check repo")
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
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
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
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
}

func TestCreateInvalidTitleNotPersisted(t *testing.T) {
	repo := newFakeRepo()
	h := New(repo)

	body := strings.NewReader(`{"title":""}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	w := httptest.NewRecorder()
	h.Create(w, req)

	tasks, _ := repo.List(context.Background())
	if len(tasks) != 0 {
		t.Errorf("в репо %d задач, want 0 — невалидная задача не должна попадать в БД", len(tasks))
	}
}

// =====================  GET BY ID  =====================

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
	assertJSONContentType(t, w)

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

func TestGetByIDNotFoundReturns404(t *testing.T) {
	h := New(newFakeRepo())
	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	h.GetByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
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
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
}

// =====================  LIST  =====================

func TestListEmptyReturnsArray(t *testing.T) {
	h := New(newFakeRepo())
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	assertJSONContentType(t, w)

	bodyStr := strings.TrimSpace(w.Body.String())
	if bodyStr == "null" {
		t.Error("body = null, want []  (пустой список — это пустой массив, не null)")
	}

	var tasks []model.Task
	if err := json.Unmarshal([]byte(bodyStr), &tasks); err != nil {
		t.Fatalf("decode: %v — body не валидный JSON-массив", err)
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
	titles := []string{"Task A", "Task B", "Task C"}
	for _, title := range titles {
		repo.Create(context.Background(), model.NewTask(title, "desc-"+title))
	}
	h := New(repo)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	assertJSONContentType(t, w)

	var tasks []model.Task
	if err := json.NewDecoder(w.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("len(tasks) = %d, want 3", len(tasks))
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

func TestListRepoErrorReturns500(t *testing.T) {
	repo := newFakeRepo()
	repo.failOn = "List"
	h := New(repo)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d при ошибке репо", w.Code, http.StatusInternalServerError)
	}
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
}

// =====================  UPDATE  =====================

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
	assertJSONContentType(t, w)

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
	if stored.Description != "New desc" {
		t.Errorf("в репозитории Description = %q, want %q", stored.Description, "New desc")
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
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
}

func TestUpdateInvalidIDReturns400(t *testing.T) {
	h := New(newFakeRepo())
	body := strings.NewReader(`{"title":"X"}`)
	req := httptest.NewRequest(http.MethodPut, "/tasks/abc", body)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
}

func TestUpdateInvalidJSONReturns400(t *testing.T) {
	repo := newFakeRepo()
	task, _ := repo.Create(context.Background(), model.NewTask("Old", ""))
	h := New(repo)

	body := strings.NewReader(`not-json`)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%d", task.ID), body)
	req.SetPathValue("id", fmt.Sprintf("%d", task.ID))
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
}

// =====================  DELETE  =====================

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
		t.Errorf("body = %q, want empty for 204 (No Content)", w.Body.String())
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
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
}

func TestDeleteInvalidIDReturns400(t *testing.T) {
	h := New(newFakeRepo())
	req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	assertJSONContentType(t, w)
	assertErrorBody(t, w)
}
