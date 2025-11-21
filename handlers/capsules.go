package handlers

import (
	"encoding/json"
	"net/http"

	"knowledge-capsule-api/middleware"
	"knowledge-capsule-api/models"
	"knowledge-capsule-api/store"
	"knowledge-capsule-api/utils"
)

var CapsuleStore = &store.CapsuleStore{FileStore: store.FileStore[models.Capsule]{FilePath: "data/capsules.json"}}

// CapsuleHandler godoc
// @Summary Get or create capsules
// @Description Get all capsules for the user or create a new one
// @Tags capsules
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param input body models.Capsule true "Capsule info (for POST)"
// @Success 200 {array} models.Capsule
// @Success 201 {object} models.Capsule
// @Failure 400 {object} map[string]interface{}
// @Router /api/capsules [get]
// @Router /api/capsules [post]
func CapsuleHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)

	switch r.Method {
	case http.MethodGet:
		capsules, _ := CapsuleStore.GetCapsulesByUser(userID)
		utils.JSONResponse(w, http.StatusOK, true, "Capsules fetched", capsules)

	case http.MethodPost:
		var req struct {
			Title     string   `json:"title"`
			Content   string   `json:"content"`
			Topic     string   `json:"topic"`
			Tags      []string `json:"tags"`
			IsPrivate bool     `json:"is_private"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		capsule, err := CapsuleStore.AddCapsule(userID, req.Title, req.Content, req.Topic, req.Tags, req.IsPrivate)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		utils.JSONResponse(w, http.StatusCreated, true, "Capsule created", capsule)

	default:
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, nil)
	}
}
