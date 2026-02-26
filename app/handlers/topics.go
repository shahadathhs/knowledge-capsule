package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"knowledge-capsule/pkg/utils"
)

// TopicHandler godoc
// @Summary Get or create topics
// @Description Get all topics (paginated) or create a new one
// @Tags topics
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20, max 100)"
// @Param input body models.Topic true "Topic info (for POST)"
// @Success 200 {object} models.PaginatedResponse "Paginated list: data, page, limit, total"
// @Success 201 {object} models.Topic
// @Failure 400 {object} map[string]interface{}
// @Router /api/topics [get]
// @Router /api/topics [post]
func TopicHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		topics, err := TopicStore.GetAllTopics()
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		page, limit := utils.ParsePagination(r)
		paged, total := utils.SlicePage(topics, page, limit)
		utils.JSONPaginatedResponse(w, http.StatusOK, "Topics fetched", paged, page, limit, total)

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
