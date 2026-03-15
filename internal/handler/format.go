package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GlebMoskalev/go-path-backend/internal/middleware"
	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/GlebMoskalev/go-path-backend/internal/utils"
	"github.com/go-chi/chi/v5"
)

type FormatHandler struct {
	formatService *service.FormatService
}

func NewFormatHandler(format *service.FormatService) *FormatHandler {
	return &FormatHandler{formatService: format}
}

func (h *FormatHandler) FormatCode(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserIDFromContext(r.Context())
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

	if req.Code == "" {
		utils.ResponseWithError(w, http.StatusBadRequest, "code is required")
		return
	}

	result, err := h.formatService.FormatCode(req.Code)
	if err != nil {
		utils.ResponseWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, struct {
		Code string `json:"code"`
	}{Code: result})
}

func (h *FormatHandler) FormatProjectCode(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectSlug := chi.URLParam(r, "projectSlug")
	stepSlug := chi.URLParam(r, "stepSlug")

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Code == "" {
		utils.ResponseWithError(w, http.StatusBadRequest, "code is required")
		return
	}

	result, err := h.formatService.FormatProjectCode(req.Code, projectSlug, stepSlug)
	if err != nil {
		utils.ResponseWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, struct {
		Code string `json:"code"`
	}{Code: result})
}
