package handlers

import (
	"net/http"

	"knowledge-capsule-api/middleware"
	"knowledge-capsule-api/utils"
)

// SearchHandler godoc
// @Summary Search capsules
// @Description Search capsules by query string
// @Tags search
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Success 200 {array} models.Capsule
// @Failure 400 {object} map[string]interface{}
// @Router /api/search [get]
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, nil)
		return
	}

	results, _ := CapsuleStore.SearchCapsules(userID, query)
	utils.JSONResponse(w, http.StatusOK, true, "Search results", results)
}
