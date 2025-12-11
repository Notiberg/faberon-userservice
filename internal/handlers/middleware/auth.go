package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/m04kA/SMC-UserService/internal/domain"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

var (
	ErrMissingUserID = errors.New("missing X-User-ID header")
	ErrInvalidUserID = errors.New("invalid user ID format")
	ErrMissingRole   = errors.New("missing X-User-Role header")
	ErrInvalidRole   = errors.New("invalid role")
)

// UserIDAuth извлекает user ID и role из заголовков X-User-ID и X-User-Role
func UserIDAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			http.Error(w, ErrMissingUserID.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, ErrInvalidUserID.Error(), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		roleStr := r.Header.Get("X-User-Role")
		if roleStr != "" {
			role := domain.Role(roleStr)
			if !role.IsValid() {
				http.Error(w, ErrInvalidRole.Error(), http.StatusBadRequest)
				return
			}
			ctx = context.WithValue(ctx, RoleKey, role)
		}

		// if roleStr == "" {
		// 	http.Error(w, ErrMissingRole.Error(), http.StatusUnauthorized)
		// 	return
		// }

		// role := domain.Role(roleStr)
		// if !role.IsValid() {
		// 	http.Error(w, ErrInvalidRole.Error(), http.StatusBadRequest)
		// 	return
		// }

		ctx = context.WithValue(ctx, UserIDKey, userID)
		// ctx = context.WithValue(ctx, RoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext извлекает user ID из контекста
func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0, ErrInvalidUserID
	}
	return userID, nil
}

// GetRoleFromContext извлекает role из контекста
func GetRoleFromContext(ctx context.Context) (domain.Role, error) {
	role, ok := ctx.Value(RoleKey).(domain.Role)
	if !ok {
		return "", ErrInvalidRole
	}
	return role, nil
}

// RequireSuperUser middleware проверяет, что пользователь имеет роль superuser
func RequireSuperUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, err := GetRoleFromContext(r.Context())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if role != domain.RoleSuperUser {
			http.Error(w, "forbidden: superuser role required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
