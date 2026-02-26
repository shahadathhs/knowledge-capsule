package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"
)

// AdminUsersHandler routes POST /api/admin/users/:id/role -> SetUserRole
func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/users/")
	if !strings.HasSuffix(path, "/role") {
		utils.ErrorResponse(w, http.StatusNotFound, nil)
		return
	}
	id := strings.TrimSuffix(path, "/role")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}
	SetUserRole(w, r)
}

// SetUserRole godoc
// @Summary Set user role (superadmin only)
// @Description Set role for a user: user, admin, or superadmin
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param input body object true "role: user|admin|superadmin"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/admin/users/{id}/role [post]
func SetUserRole(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodPost) {
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/admin/users/")
	id := strings.TrimSuffix(path, "/role")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if r.Body == nil {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	role := strings.TrimSpace(req.Role)
	if role == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}

	if err := UserStore.UpdateUserRole(id, role); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, true, "Role updated", map[string]string{"user_id": id, "role": role})
}

// ListAdmins godoc
// @Summary List admins (superadmin only)
// @Description List users with role admin or superadmin
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} models.PaginatedResponse
// @Failure 403 {object} map[string]interface{}
// @Router /api/admin/admins [get]
func ListAdmins(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}
	page, limit := utils.ParsePagination(r)

	admins, total, err := UserStore.ListAdmins(page, limit)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSONPaginatedResponse(w, http.StatusOK, "Admins fetched", admins, page, limit, total)
}

// GlobalSearchResult is the response for admin global search.
type GlobalSearchResult struct {
	Users    []models.User    `json:"users"`
	Topics   []models.Topic   `json:"topics"`
	Capsules []models.Capsule `json:"capsules"`
}

// GlobalSearch godoc
// @Summary Global search (admin only)
// @Description Search across users, topics, and capsules
// @Tags admin
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param limit query int false "Max results per type (default 10)"
// @Success 200 {object} GlobalSearchResult
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/admin/search [get]
func GlobalSearch(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 50 {
			limit = v
		}
	}

	users, _ := UserStore.SearchUsers(q, limit)
	topics, _ := TopicStore.SearchTopics(q, limit)
	capsules, _ := CapsuleStore.SearchAllCapsules(q, limit)

	_ = r.Context().Value(middleware.UserContextKey) // ensure auth

	result := GlobalSearchResult{
		Users:    users,
		Topics:   topics,
		Capsules: capsules,
	}
	utils.JSONResponse(w, http.StatusOK, true, "Search results", result)
}
