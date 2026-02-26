package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/logger"
	"knowledge-capsule/pkg/utils"
)

// GetTopics godoc
// @Summary Get topics
// @Description Get all topics (paginated, filterable)
// @Tags topics
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20, max 100)"
// @Param q query string false "Search in name or description"
// @Success 200 {object} models.PaginatedResponse "Paginated list: data, page, limit, total"
// @Failure 400 {object} map[string]interface{}
// @Router /api/topics [get]
func GetTopics(w http.ResponseWriter, r *http.Request) {
	var filters *models.TopicFilters
	if q := r.URL.Query().Get("q"); q != "" {
		filters = &models.TopicFilters{Q: q}
	}

	topics, err := TopicStore.GetAllTopics(filters)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	page, limit := utils.ParsePagination(r)
	paged, total := utils.SlicePage(topics, page, limit)
	logger.LogEvent(logger.EventTopic, r, slog.String("action", "list"), slog.Int("count", len(paged)), slog.Int("total", total))
	utils.JSONPaginatedResponse(w, http.StatusOK, "Topics fetched", paged, page, limit, total)
}

// CreateTopic godoc
// @Summary Create topic
// @Description Create a new topic
// @Tags topics
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param input body models.TopicInput true "Topic data"
// @Success 201 {object} models.Topic
// @Failure 400 {object} map[string]interface{}
// @Router /api/topics [post]
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req models.TopicInput
	json.NewDecoder(r.Body).Decode(&req)
	topic, err := TopicStore.AddTopic(req.Name, req.Description)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}
	logger.LogEvent(logger.EventTopic, r, slog.String("action", "create"), slog.String("topic_id", topic.ID), slog.String("name", req.Name))
	utils.JSONResponse(w, http.StatusCreated, true, "Topic created", topic)
}

// TopicHandler routes GET/POST to GetTopics or CreateTopic.
func TopicHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTopics(w, r)
	case http.MethodPost:
		CreateTopic(w, r)
	default:
		utils.ErrorResponse(w, r, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}
