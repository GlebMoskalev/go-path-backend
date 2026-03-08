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
