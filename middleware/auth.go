package middleware

import (
	"context"
	"net/http"
	"strings"

	"knowledge-capsule-api/utils"
)

type contextKey string

const UserContextKey = contextKey("user_id")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
