package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/app/models"
	"knowledge-capsule/app/store"
	"knowledge-capsule/pkg/utils"
)

const maxTitleLen = 500

var CapsuleStore = &store.CapsuleStore{FileStore: store.FileStore[models.Capsule]{FilePath: "data/capsules.json"}}

// CapsuleHandler godoc
// @Summary Get or create capsules
// @Description Get all capsules for the user (paginated) or create a new one
// @Tags capsules
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20, max 100)"
// @Param input body models.Capsule true "Capsule info (for POST)"
// @Success 200 {object} models.PaginatedResponse "Paginated list: data, page, limit, total"
// @Success 201 {object} models.Capsule
// @Failure 400 {object} map[string]interface{}
// @Router /api/capsules [get]
// @Router /api/capsules [post]
func CapsuleHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)

	switch r.Method {
	case http.MethodGet:
		capsules, err := CapsuleStore.GetCapsulesByUser(userID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		page, limit := utils.ParsePagination(r)
		paged, total := utils.SlicePage(capsules, page, limit)
		utils.JSONPaginatedResponse(w, http.StatusOK, "Capsules fetched", paged, page, limit, total)

	case http.MethodPost:
		var req struct {
			Title     string   `json:"title"`
			Content   string   `json:"content"`
			Topic     string   `json:"topic"`
			Tags      []string `json:"tags"`
			IsPrivate bool     `json:"is_private"`
		}
		if r.Body == nil {
			utils.ErrorResponse(w, http.StatusBadRequest, errors.New("empty request body"))
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		title := strings.TrimSpace(req.Title)
		if title == "" {
			utils.ErrorResponse(w, http.StatusBadRequest, &utils.ValidationError{Field: "title", Message: "cannot be empty"})
			return
		}
		if len(title) > maxTitleLen {
			utils.ErrorResponse(w, http.StatusBadRequest, &utils.ValidationError{Field: "title", Message: "exceeds maximum length"})
			return
		}

		capsule, err := CapsuleStore.AddCapsule(userID, title, req.Content, req.Topic, req.Tags, req.IsPrivate)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		utils.JSONResponse(w, http.StatusCreated, true, "Capsule created", capsule)

	default:
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, nil)
	}
}
