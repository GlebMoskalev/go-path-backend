package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GlebMoskalev/go-path-backend/internal/middleware"
	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/GlebMoskalev/go-path-backend/internal/utils"
)

type FormatHandler struct {
	formatService *service.FormatService
}

func NewFormatHandler(formatService *service.FormatService) *FormatHandler {
	return &FormatHandler{formatService: formatService}
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

	if len(req.Code) == 0 {
		utils.ResponseWithError(w, http.StatusBadRequest, "code is required")
		return
	}

	result, err := h.formatService.FormatCode(req.Code)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	formattedCode := struct {
		Code string `json:"code"`
	}{
		Code: result,
	}
	utils.ResponseWithJSON(w, http.StatusOK, formattedCode)
}
