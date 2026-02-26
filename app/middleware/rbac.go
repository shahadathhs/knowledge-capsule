package middleware

import (
	"net/http"

	"knowledge-capsule/app/models"
)

// RequireAdmin wraps a handler and returns 403 if the user's role is not admin or superadmin.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(RoleContextKey).(string)
		if role != models.RoleAdmin && role != models.RoleSuperAdmin {
			http.Error(w, "forbidden: admin access required", http.StatusForbidden)
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
			http.Error(w, "forbidden: superadmin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
