package middleware

import (
	"context"
	"net/http"
	"strings"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"
)

type contextKey string

const (
	UserContextKey = contextKey("user_id")
	RoleContextKey = contextKey("user_role")
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		var tokenString string

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = r.URL.Query().Get("token")
		}

		if tokenString == "" {
			http.Error(w, "missing or invalid Authorization header or token", http.StatusUnauthorized)
			return
		}
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		role := claims.Role
		if role == "" {
			role = models.RoleUser
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserContextKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleContextKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
