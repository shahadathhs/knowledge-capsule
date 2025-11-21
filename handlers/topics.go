package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"knowledge-capsule-api/models"
	"knowledge-capsule-api/store"
	"knowledge-capsule-api/utils"
)

var TopicStore = &store.TopicStore{FileStore: store.FileStore[models.Topic]{FilePath: "data/topics.json"}}

// TopicHandler godoc
// @Summary Get or create topics
// @Description Get all topics or create a new one
// @Tags topics
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param input body models.Topic true "Topic info (for POST)"
// @Success 200 {array} models.Topic
// @Success 201 {object} models.Topic
// @Failure 400 {object} map[string]interface{}
// @Router /api/topics [get]
// @Router /api/topics [post]
func TopicHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		topics, _ := TopicStore.GetAllTopics()
		utils.JSONResponse(w, http.StatusOK, true, "Topics fetched", topics)

	case http.MethodPost:
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		topic, err := TopicStore.AddTopic(req.Name, req.Description)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		utils.JSONResponse(w, http.StatusCreated, true, "Topic created", topic)

	default:
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}
