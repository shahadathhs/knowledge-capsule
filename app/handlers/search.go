package handlers

import (
	"net/http"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/pkg/utils"
)

// SearchHandler godoc
// @Summary Search capsules
// @Description Search capsules by query string (paginated)
// @Tags search
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20, max 100)"
// @Success 200 {object} models.PaginatedResponse "Paginated list: data, page, limit, total"
// @Failure 400 {object} map[string]interface{}
// @Router /api/search [get]
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}

	results, err := CapsuleStore.SearchCapsules(userID, query)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	page, limit := utils.ParsePagination(r)
	paged, total := utils.SlicePage(results, page, limit)
	utils.JSONPaginatedResponse(w, http.StatusOK, "Search results", paged, page, limit, total)
}
