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

type ProjectHandler struct {
	projectService    *service.ProjectService
	submissionService *service.SubmissionService
}

func NewProjectHandler(
	projectService *service.ProjectService,
	submissionService *service.SubmissionService,
) *ProjectHandler {
	return &ProjectHandler{
		projectService:    projectService,
		submissionService: submissionService,
	}
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	userID := middleware.OptionalUserID(r.Context())
	projects := h.projectService.ListProjects(r.Context(), userID)
	utils.ResponseWithJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "projectSlug")
	userID := middleware.OptionalUserID(r.Context())

	project, err := h.projectService.GetProject(r.Context(), slug, userID)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "project not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) GetStep(w http.ResponseWriter, r *http.Request) {
	projectSlug := chi.URLParam(r, "projectSlug")
	stepSlug := chi.URLParam(r, "stepSlug")
	userID := middleware.OptionalUserID(r.Context())

	step, err := h.projectService.GetStep(r.Context(), projectSlug, stepSlug, userID)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) || errors.Is(err, service.ErrProjectStepNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "step not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, step)
}

func (h *ProjectHandler) Submit(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.submissionService.SubmitProject(r.Context(), userID, projectSlug, stepSlug, req.Code)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) || errors.Is(err, service.ErrProjectStepNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "step not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, result)
}
