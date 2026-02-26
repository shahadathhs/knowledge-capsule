package middleware

import (
	"errors"
	"net/http"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"
)

// RequireAdmin wraps a handler and returns 403 if the user's role is not admin or superadmin.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(RoleContextKey).(string)
		if role != models.RoleAdmin && role != models.RoleSuperAdmin {
			utils.ErrorResponse(w, r, http.StatusForbidden, errors.New("admin access required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireSuperAdmin wraps a handler and returns 403 if the user's role is not superadmin.
func RequireSuperAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(RoleContextKey).(string)
		if role != models.RoleSuperAdmin {
			utils.ErrorResponse(w, r, http.StatusForbidden, errors.New("superadmin access required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
