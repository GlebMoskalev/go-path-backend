package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"taskmanager/model"
)

type TaskRepository interface {
	Create(ctx context.Context, task model.Task) (model.Task, error)
	GetByID(ctx context.Context, id int64) (model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
	Update(ctx context.Context, task model.Task) (model.Task, error)
	Delete(ctx context.Context, id int64) error
}

type TaskHandler struct {
	repo TaskRepository
}

func New(repo TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := task.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	task.Status = model.StatusPending
	created, err := h.repo.Create(r.Context(), task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create task")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}
	if tasks == nil {
		tasks = make([]model.Task, 0)
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	task.ID = id
	updated, err := h.repo.Update(r.Context(), task)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
