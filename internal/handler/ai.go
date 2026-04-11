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

type AIHandler struct {
	aiService *service.AIService
}

func NewAIHandler(aiService *service.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

func (h *AIHandler) AnalyzePassedTask(w http.ResponseWriter, r *http.Request) {
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

	recommendation, err := h.aiService.AnalyzePassedCodeTask(r.Context(), chapterSlug, taskSlug, req.Code, userID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotPassed) {
			utils.ResponseWithError(w, http.StatusForbidden, "task not passed yet")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := struct {
		Recommendation string `json:"recommendation"`
	}{
		Recommendation: recommendation,
	}

	utils.ResponseWithJSON(w, http.StatusOK, resp)
}

func (h *AIHandler) AnalyzePassedProject(w http.ResponseWriter, r *http.Request) {
	projectSlug := chi.URLParam(r, "projectSlug")
	stepSlug := chi.URLParam(r, "stepSlug")
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

	recommendation, err := h.aiService.AnalyzePassedCodeProject(r.Context(), projectSlug, stepSlug, req.Code, userID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotPassed) {
			utils.ResponseWithError(w, http.StatusForbidden, "step not passed yet")
			return
		}
		if errors.Is(err, service.ErrProjectNotFound) || errors.Is(err, service.ErrProjectStepNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "step not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, struct {
		Recommendation string `json:"recommendation"`
	}{Recommendation: recommendation})
}

func (h *AIHandler) AnalyzeErrorTask(w http.ResponseWriter, r *http.Request) {
	chapterSlug := chi.URLParam(r, "chapterSlug")
	taskSlug := chi.URLParam(r, "taskSlug")
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Code  string `json:"code"`
		Error string `json:"error"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Code) == 0 {
		utils.ResponseWithError(w, http.StatusBadRequest, "code is required")
		return
	}

	if len(req.Error) == 0 {
		utils.ResponseWithError(w, http.StatusBadRequest, "error is required")
		return
	}

	if len(req.Code) > 10240 {
		utils.ResponseWithError(w, http.StatusBadRequest, "code too large")
		return
	}

	if len(req.Error) > 10240 {
		utils.ResponseWithError(w, http.StatusBadRequest, "error output too large")
		return
	}

	analysis, err := h.aiService.AnalyzeErrorTask(r.Context(), chapterSlug, taskSlug, req.Code, req.Error, userID)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := struct {
		Analysis string `json:"analysis"`
	}{
		Analysis: analysis,
	}

	utils.ResponseWithJSON(w, http.StatusOK, resp)
}

func (h *AIHandler) AnalyzeErrorProject(w http.ResponseWriter, r *http.Request) {
	projectSlug := chi.URLParam(r, "projectSlug")
	stepSlug := chi.URLParam(r, "stepSlug")
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Code  string `json:"code"`
		Error string `json:"error"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Code) == 0 {
		utils.ResponseWithError(w, http.StatusBadRequest, "code is required")
		return
	}

	if len(req.Error) == 0 {
		utils.ResponseWithError(w, http.StatusBadRequest, "error is required")
		return
	}

	if len(req.Code) > 10240 {
		utils.ResponseWithError(w, http.StatusBadRequest, "code too large")
		return
	}

	if len(req.Error) > 10240 {
		utils.ResponseWithError(w, http.StatusBadRequest, "error output too large")
		return
	}

	analysis, err := h.aiService.AnalyzeErrorProject(r.Context(), projectSlug, stepSlug, req.Code, req.Error, userID)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) || errors.Is(err, service.ErrProjectStepNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "step not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, struct {
		Analysis string `json:"analysis"`
	}{Analysis: analysis})
}
