package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"todoapi/model"
	"todoapi/service"
	"todoapi/storage"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(svc *service.TodoService) *TodoHandler {
	return &TodoHandler{service: svc}
}

func (h *TodoHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", h.GetAll)
	mux.HandleFunc("POST /todos", h.Create)
	mux.HandleFunc("GET /todos/{id}", h.GetByID)
	mux.HandleFunc("PUT /todos/{id}", h.Update)
	mux.HandleFunc("DELETE /todos/{id}", h.Delete)
	return mux
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo := h.service.Create(req)
	writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	todos := h.service.GetAll()
	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	todo, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todo, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
