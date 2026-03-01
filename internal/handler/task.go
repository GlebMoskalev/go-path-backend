package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GlebMoskalev/go-path-backend/internal/middleware"
	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/GlebMoskalev/go-path-backend/internal/utils"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService       *service.TaskService
	submissionService *service.SubmissionService
}

func NewTaskHandler(
	taskService *service.TaskService,
	submissionService *service.SubmissionService,
) *TaskHandler {
	return &TaskHandler{
		taskService:       taskService,
		submissionService: submissionService,
	}
}
func (h *TaskHandler) ListChapters(w http.ResponseWriter, r *http.Request) {
	userID := middleware.OptionalUserID(r.Context())
	chapters := h.taskService.ListChapters(r.Context(), userID)
	utils.ResponseWithJSON(w, http.StatusOK, chapters)
}

func (h *TaskHandler) GetChapter(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "chapterSlug")
	userID := middleware.OptionalUserID(r.Context())

	chapter, err := h.taskService.GetChapter(r.Context(), slug, userID)
	if err != nil {
		if errors.Is(err, service.ErrTaskChapterNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "chapter not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, chapter)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	chapterSlug := chi.URLParam(r, "chapterSlug")
	taskSlug := chi.URLParam(r, "taskSlug")
	userID := middleware.OptionalUserID(r.Context())

	task, err := h.taskService.GetTask(r.Context(), chapterSlug, taskSlug, userID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) || errors.Is(err, service.ErrTaskChapterNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "task not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Submit(w http.ResponseWriter, r *http.Request) {
	chapterSlug := chi.URLParam(r, "chapterSlug")
	taskSlug := chi.URLParam(r, "taskSlug")
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Code) == 0 {
		utils.ResponseWithError(w, http.StatusBadRequest, "code is required")
		return
	}

	if len(req.Code) > 10240 {
		utils.ResponseWithError(w, http.StatusBadRequest, "code too large")
		return
	}

	result, err := h.submissionService.Submit(r.Context(), userID, chapterSlug, taskSlug, req.Code)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) || errors.Is(err, service.ErrTaskChapterNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "task not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, result)
}
