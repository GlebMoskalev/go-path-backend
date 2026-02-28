package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GlebMoskalev/go-path-backend/internal/middleware"
	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/GlebMoskalev/go-path-backend/internal/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			utils.ResponseWithError(w, http.StatusNotFound, "user not found")
			return
		}
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.ResponseWithError(w, http.StatusBadRequest, "name is required")
		return
	}

	if err := h.userService.UpdateProfile(r.Context(), userID, req.Name, req.Picture); err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "profile updated successfully"})
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.userService.DeleteAccount(r.Context(), userID); err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, "failed to delete account")
		return
	}

	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"message": "account deleted successfully"})
}
