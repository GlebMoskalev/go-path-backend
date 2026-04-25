package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"taskmanager/model"
)

func TestHealthCheck(t *testing.T) {
	srv := SetupTestServer(t)

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status field = %q, want %q", body["status"], "ok")
	}
}

func TestCreateAndGetTask(t *testing.T) {
	srv := SetupTestServer(t)

	resp, err := http.Post(srv.URL+"/tasks", "application/json",
		strings.NewReader(`{"title":"Integration task","description":"Some description"}`))
	if err != nil {
		t.Fatalf("POST /tasks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("create status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	var created model.Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("created.ID = 0, want non-zero")
	}

	resp2, err := http.Get(fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID))
	if err != nil {
		t.Fatalf("GET /tasks/%d: %v", created.ID, err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("get status = %d, want %d", resp2.StatusCode, http.StatusOK)
	}

	var got model.Task
	if err := json.NewDecoder(resp2.Body).Decode(&got); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if got.Title != "Integration task" {
		t.Errorf("Title = %q, want %q", got.Title, "Integration task")
	}
	if got.Status != model.StatusPending {
		t.Errorf("Status = %q, want %q", got.Status, model.StatusPending)
	}
}

func TestListTasks(t *testing.T) {
	srv := SetupTestServer(t)

	for _, title := range []string{"Task A", "Task B", "Task C"} {
		resp, err := http.Post(srv.URL+"/tasks", "application/json",
			strings.NewReader(fmt.Sprintf(`{"title":%q}`, title)))
		if err != nil {
			t.Fatalf("create %q: %v", title, err)
		}
		resp.Body.Close()
	}

	resp, err := http.Get(srv.URL + "/tasks")
	if err != nil {
		t.Fatalf("GET /tasks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var tasks []model.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("len(tasks) = %d, want 3", len(tasks))
	}
}

func TestUpdateTask(t *testing.T) {
	srv := SetupTestServer(t)

	resp, err := http.Post(srv.URL+"/tasks", "application/json",
		strings.NewReader(`{"title":"Original","description":""}`))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	var created model.Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	resp.Body.Close()

	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID),
		strings.NewReader(`{"title":"Updated","description":"New desc","status":"in_progress"}`))
	if err != nil {
		t.Fatalf("build PUT: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /tasks/%d: %v", created.ID, err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("update status = %d, want %d", resp2.StatusCode, http.StatusOK)
	}

	var updated model.Task
	if err := json.NewDecoder(resp2.Body).Decode(&updated); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	if updated.Title != "Updated" {
		t.Errorf("Title = %q, want %q", updated.Title, "Updated")
	}
	if updated.Status != model.StatusInProgress {
		t.Errorf("Status = %q, want %q", updated.Status, model.StatusInProgress)
	}
}

func TestDeleteTask(t *testing.T) {
	srv := SetupTestServer(t)

	resp, err := http.Post(srv.URL+"/tasks", "application/json",
		strings.NewReader(`{"title":"To delete"}`))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	var created model.Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	resp.Body.Close()

	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID), nil)
	if err != nil {
		t.Fatalf("build DELETE: %v", err)
	}
	resp2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE: %v", err)
	}
	resp2.Body.Close()

	if resp2.StatusCode != http.StatusNoContent {
		t.Errorf("delete status = %d, want %d", resp2.StatusCode, http.StatusNoContent)
	}

	resp3, err := http.Get(fmt.Sprintf("%s/tasks/%d", srv.URL, created.ID))
	if err != nil {
		t.Fatalf("GET after delete: %v", err)
	}
	defer resp3.Body.Close()
	if resp3.StatusCode != http.StatusNotFound {
		t.Errorf("get after delete = %d, want %d", resp3.StatusCode, http.StatusNotFound)
	}
}

func TestValidation(t *testing.T) {
	srv := SetupTestServer(t)

	resp, err := http.Post(srv.URL+"/tasks", "application/json",
		strings.NewReader(`{"title":""}`))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d for empty title", resp.StatusCode, http.StatusBadRequest)
	}
}
