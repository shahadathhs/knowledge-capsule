package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/pkg/utils"
)

// TopicByIDHandler godoc
// @Summary Get, update, or delete a topic by ID
// @Description Get, update, or delete a single topic
// @Tags topics
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Topic ID"
// @Param input body models.Topic true "Updated topic fields (for PUT)"
// @Success 200 {object} models.Topic "Topic fetched or updated"
// @Success 404 {object} map[string]interface{} "Topic not found"
// @Router /api/topics/{id} [get]
// @Router /api/topics/{id} [put]
// @Router /api/topics/{id} [delete]
func TopicByIDHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.UserContextKey).(string) // auth required
	id := strings.TrimPrefix(r.URL.Path, "/api/topics/")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("missing topic id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		topic, err := TopicStore.FindByID(id)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, true, "Topic fetched", topic)

	case http.MethodPut:
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if r.Body == nil {
			utils.ErrorResponse(w, http.StatusBadRequest, errors.New("empty request body"))
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		name := strings.TrimSpace(req.Name)
		if name == "" {
			utils.ErrorResponse(w, http.StatusBadRequest, &utils.ValidationError{Field: "name", Message: "cannot be empty"})
			return
		}

		topic, err := TopicStore.UpdateTopic(id, name, req.Description)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, true, "Topic updated", topic)

	case http.MethodDelete:
		if err := TopicStore.DeleteTopic(id); err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err)
			return
		}
		utils.JSONResponse(w, http.StatusOK, true, "Topic deleted", nil)

	default:
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, nil)
	}
}
