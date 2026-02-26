package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/utils"
)

// CapsuleByIDHandler godoc
// @Summary Get, update, or delete a capsule by ID
// @Description Get, update, or delete a single capsule (user must own it)
// @Tags capsules
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Capsule ID"
// @Param input body models.Capsule true "Updated capsule fields (for PUT)"
// @Success 200 {object} models.Capsule "Capsule fetched or updated"
// @Success 404 {object} map[string]interface{} "Capsule not found"
// @Router /api/capsules/{id} [get]
// @Router /api/capsules/{id} [put]
// @Router /api/capsules/{id} [delete]
func CapsuleByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)
	id := strings.TrimPrefix(r.URL.Path, "/api/capsules/")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("missing capsule id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		capsule, err := CapsuleStore.FindByID(id)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err)
			return
		}
		if capsule.UserID != userID {
			utils.ErrorResponse(w, http.StatusNotFound, errors.New("capsule not found"))
			return
		}
		utils.JSONResponse(w, http.StatusOK, true, "Capsule fetched", capsule)

	case http.MethodPut:
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

		updated := models.Capsule{
			Title:     title,
			Content:   req.Content,
			Topic:     req.Topic,
			Tags:      models.Tags(req.Tags),
			IsPrivate: req.IsPrivate,
		}
		capsule, err := CapsuleStore.UpdateCapsule(id, userID, updated)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, true, "Capsule updated", capsule)

	case http.MethodDelete:
		if err := CapsuleStore.DeleteCapsule(id, userID); err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, true, "Capsule deleted", nil)

	default:
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, nil)
	}
}
