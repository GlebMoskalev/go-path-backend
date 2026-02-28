package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/GlebMoskalev/go-path-backend/internal/service"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "user_id"

type AuthMiddleware struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthMiddleware(authService *service.AuthService, userService *service.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		userService: userService,
	}
}

// Authenticate проверяет JWT токен и добавляет user id в контекст
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		userID, tokenVersion, err := m.authService.ValidateAccessToken(r.Context(), token)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		user, err := m.userService.GetByID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if tokenVersion != user.TokenVersion {
			http.Error(w, "token has been invalidated", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(userIDKey).(uuid.UUID)
	return value, ok
}
