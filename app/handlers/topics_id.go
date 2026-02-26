package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"knowledge-capsule/app/middleware"
	"knowledge-capsule/app/models"
	"knowledge-capsule/pkg/logger"
	"knowledge-capsule/pkg/utils"
)

// GetTopicByID godoc
// @Summary Get topic by ID
// @Description Get a single topic
// @Tags topics
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Topic ID"
// @Success 200 {object} models.Topic
// @Failure 404 {object} map[string]interface{}
// @Router /api/topics/{id} [get]
func GetTopicByID(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.UserContextKey).(string) // auth required
	id := strings.TrimPrefix(r.URL.Path, "/api/topics/")
	if id == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("missing topic id"))
		return
	}
	topic, err := TopicStore.FindByID(id)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusNotFound, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, true, "Topic fetched", topic)
}

// UpdateTopicByID godoc
// @Summary Update topic by ID
// @Description Update a topic
// @Tags topics
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Topic ID"
// @Param input body models.TopicInput true "Updated topic fields"
// @Success 200 {object} models.Topic
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/topics/{id} [put]
func UpdateTopicByID(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.UserContextKey).(string) // auth required
	id := strings.TrimPrefix(r.URL.Path, "/api/topics/")
	if id == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("missing topic id"))
		return
	}
	var req models.TopicInput
	if r.Body == nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("empty request body"))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, &utils.ValidationError{Field: "name", Message: "cannot be empty"})
		return
	}

	topic, err := TopicStore.UpdateTopic(id, name, req.Description)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusNotFound, err)
		return
	}
	logger.LogEvent(logger.EventTopic, r, slog.String("action", "update"), slog.String("topic_id", id))
	utils.JSONResponse(w, http.StatusOK, true, "Topic updated", topic)
}

// DeleteTopicByID godoc
// @Summary Delete topic by ID
// @Description Delete a topic
// @Tags topics
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Topic ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/topics/{id} [delete]
func DeleteTopicByID(w http.ResponseWriter, r *http.Request) {
	_ = r.Context().Value(middleware.UserContextKey).(string) // auth required
	id := strings.TrimPrefix(r.URL.Path, "/api/topics/")
	if id == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("missing topic id"))
		return
	}
	if err := TopicStore.DeleteTopic(id); err != nil {
		utils.ErrorResponse(w, r, http.StatusNotFound, err)
		return
	}
	logger.LogEvent(logger.EventTopic, r, slog.String("action", "delete"), slog.String("topic_id", id))
	utils.JSONResponse(w, http.StatusOK, true, "Topic deleted", nil)
}

// TopicByIDHandler routes GET/PUT/DELETE to the appropriate handler.
func TopicByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTopicByID(w, r)
	case http.MethodPut:
		UpdateTopicByID(w, r)
	case http.MethodDelete:
		DeleteTopicByID(w, r)
	default:
		utils.ErrorResponse(w, r, http.StatusMethodNotAllowed, nil)
	}
}
