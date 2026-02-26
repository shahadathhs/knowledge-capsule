package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"
)

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the authenticated user's profile
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]interface{}
// @Router /api/users/me [get]
func GetProfile(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}
	userID := r.Context().Value(middleware.UserContextKey).(string)
	user, err := UserStore.FindByID(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, true, "Profile fetched", user)
}

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Update name and/or avatar_url
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param input body object true "name, avatar_url (optional)"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Router /api/users/me [patch]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodPatch) {
		return
	}
	userID := r.Context().Value(middleware.UserContextKey).(string)

	var req struct {
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}
	if r.Body == nil {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	name := strings.TrimSpace(req.Name)
	if err := UserStore.UpdateProfile(userID, name, req.AvatarURL); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	user, _ := UserStore.FindByID(userID)
	utils.JSONResponse(w, http.StatusOK, true, "Profile updated", user)
}

// ListUsers godoc
// @Summary List users (admin only)
// @Description Paginated list with optional search and role filter
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param q query string false "Search by name or email"
// @Param role query string false "Filter by role (user, admin, superadmin)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} models.PaginatedResponse
// @Failure 403 {object} map[string]interface{}
// @Router /api/users [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}
	q := r.URL.Query().Get("q")
	role := r.URL.Query().Get("role")
	page, limit := utils.ParsePagination(r)

	users, total, err := UserStore.ListUsers(q, role, page, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSONPaginatedResponse(w, http.StatusOK, "Users fetched", users, page, limit, total)
}

// GetUserByID godoc
// @Summary Get user by ID (admin only)
// @Description Get any user's profile by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/users/{id} [get]
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}
	// Admin check is done by RequireAdmin middleware on this route
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}

	user, err := UserStore.FindByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, true, "User fetched", user)
}

// UserHandler routes /api/users/me (GET, PATCH) and /api/users/:id (GET, admin only).
func UserHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/users")
	path = strings.TrimPrefix(path, "/")

	if path == "me" {
		switch r.Method {
		case http.MethodGet:
			GetProfile(w, r)
		case http.MethodPatch:
			UpdateProfile(w, r)
		default:
			utils.ErrorResponse(w, http.StatusMethodNotAllowed, nil)
		}
		return
	}

	// /api/users/:id - admin only
	if path != "" {
		role, _ := r.Context().Value(middleware.RoleContextKey).(string)
		if role != models.RoleAdmin && role != models.RoleSuperAdmin {
			http.Error(w, "forbidden: admin access required", http.StatusForbidden)
			return
		}
		GetUserByID(w, r)
		return
	}

	utils.ErrorResponse(w, http.StatusNotFound, nil)
}
