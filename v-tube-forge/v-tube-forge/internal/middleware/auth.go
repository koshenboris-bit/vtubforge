package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend-service/internal/domain"
	"backend-service/internal/repository"
	"backend-service/internal/service"
)

type contextKey string

const (
	ContextUserID  contextKey = "userID"
	ContextUserRole contextKey = "userRole"
)

func Auth(secret string, users *repository.UserRepository) func(http.Handler) http.Handler {
	_ = users
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			userID, role, err := service.ParseAccessToken(secret, token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserID, userID)
			ctx = context.WithValue(ctx, ContextUserRole, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got, _ := r.Context().Value(ContextUserRole).(string)
			if got != role {
				writeError(w, http.StatusForbidden, domain.ErrForbidden.Error())
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserID(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(ContextUserID).(int64)
	return v, ok
}

func UserRole(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ContextUserRole).(string)
	return v, ok
}
