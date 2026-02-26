package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/contextkeys"
	"knowledge-capsule/pkg/utils"
)

// Re-export for handler convenience
var (
	UserContextKey = contextkeys.UserContextKey
	RoleContextKey = contextkeys.RoleContextKey
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
			utils.ErrorResponse(w, r, http.StatusUnauthorized, errors.New("missing or invalid Authorization header or token"))
			return
		}
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusUnauthorized, err)
			return
		}

		role := claims.Role
		if role == "" {
			role = models.RoleUser
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextkeys.UserContextKey, claims.UserID)
		ctx = context.WithValue(ctx, contextkeys.RoleContextKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
