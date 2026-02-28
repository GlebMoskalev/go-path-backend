package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GlebMoskalev/go-path-backend/internal/middleware"
	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/GlebMoskalev/go-path-backend/internal/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url, err := h.authService.GetGoogleLoginURL(r.Context())
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to generate login url")
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if code == "" || state == "" {
		utils.ResponseWithError(w, http.StatusBadRequest, "missing code or state")
		return
	}

	tokenPair, err := h.authService.HandleGoogleCallback(r.Context(), code, state)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to authenticate")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, tokenPair)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokenPair, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		utils.ResponseWithError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, tokenPair)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.authService.Logout(r.Context(), userID); err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}
