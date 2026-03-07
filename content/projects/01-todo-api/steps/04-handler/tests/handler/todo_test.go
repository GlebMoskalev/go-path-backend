package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todoapi/model"
	"todoapi/service"
	"todoapi/storage"
)

func newTestServer() *httptest.Server {
	store := storage.NewTodoStorage()
	svc := service.NewTodoService(store)
	h := NewTodoHandler(svc)
	return httptest.NewServer(h.Routes())
}

func TestHandlerCreate(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"title":"Buy milk","description":"Go to store"}`
	resp, err := http.Post(srv.URL+"/todos", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("POST /todos: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	var todo model.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if todo.ID == 0 {
		t.Error("ID should not be 0")
	}
	if todo.Title != "Buy milk" {
		t.Errorf("Title = %q, want %q", todo.Title, "Buy milk")
	}
	if todo.Description != "Go to store" {
		t.Errorf("Description = %q, want %q", todo.Description, "Go to store")
	}
}

func TestHandlerCreateEmptyTitle(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"title":"","description":"test"}`
	resp, err := http.Post(srv.URL+"/todos", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("POST /todos: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestHandlerGetAll(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	for _, title := range []string{"First", "Second"} {
		body := `{"title":"` + title + `"}`
		resp, err := http.Post(srv.URL+"/todos", "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("POST /todos: %v", err)
		}
		resp.Body.Close()
	}

	resp, err := http.Get(srv.URL + "/todos")
	if err != nil {
		t.Fatalf("GET /todos: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var todos []model.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("len = %d, want 2", len(todos))
	}
}

func TestHandlerGetAllEmpty(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/todos")
	if err != nil {
		t.Fatalf("GET /todos: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var todos []model.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("len = %d, want 0", len(todos))
	}
}

func TestHandlerGetByID(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"title":"Test"}`
	postResp, err := http.Post(srv.URL+"/todos", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("POST /todos: %v", err)
	}
	postResp.Body.Close()

	resp, err := http.Get(srv.URL + "/todos/1")
	if err != nil {
		t.Fatalf("GET /todos/1: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var todo model.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if todo.Title != "Test" {
		t.Errorf("Title = %q, want %q", todo.Title, "Test")
	}
}

func TestHandlerGetByIDNotFound(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/todos/999")
	if err != nil {
		t.Fatalf("GET /todos/999: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestHandlerUpdate(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	createBody := `{"title":"Original"}`
	postResp, err := http.Post(srv.URL+"/todos", "application/json", bytes.NewBufferString(createBody))
	if err != nil {
		t.Fatalf("POST /todos: %v", err)
	}
	postResp.Body.Close()

	updateBody := `{"title":"Updated","done":true}`
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/todos/1", bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /todos/1: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var todo model.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if todo.Title != "Updated" {
		t.Errorf("Title = %q, want %q", todo.Title, "Updated")
	}
	if !todo.Done {
		t.Error("Done = false, want true")
	}
}

func TestHandlerUpdateNotFound(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"title":"New"}`
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/todos/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /todos/999: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestHandlerDelete(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	body := `{"title":"Test"}`
	postResp, err := http.Post(srv.URL+"/todos", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("POST /todos: %v", err)
	}
	postResp.Body.Close()

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/todos/1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /todos/1: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}

	getResp, err := http.Get(srv.URL + "/todos/1")
	if err != nil {
		t.Fatalf("GET /todos/1: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("after delete: status = %d, want %d", getResp.StatusCode, http.StatusNotFound)
	}
}

func TestHandlerDeleteNotFound(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/todos/999", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /todos/999: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}
