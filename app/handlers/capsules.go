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

const maxTitleLen = 500

// GetCapsules godoc
// @Summary Get capsules
// @Description Get all capsules for the user (paginated, filterable)
// @Tags capsules
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20, max 100)"
// @Param topic query string false "Filter by topic"
// @Param tags query string false "Filter by tags (comma-separated)"
// @Param q query string false "Search in title, content, tags"
// @Param is_private query bool false "Filter by is_private"
// @Success 200 {object} models.PaginatedResponse "Paginated list: data, page, limit, total"
// @Failure 400 {object} map[string]interface{}
// @Router /api/capsules [get]
func GetCapsules(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)

	var filters *models.CapsuleFilters
	if topic := r.URL.Query().Get("topic"); topic != "" ||
		r.URL.Query().Get("tags") != "" ||
		r.URL.Query().Get("q") != "" ||
		r.URL.Query().Get("is_private") != "" {
		filters = &models.CapsuleFilters{
			Topic: r.URL.Query().Get("topic"),
			Q:     r.URL.Query().Get("q"),
		}
		if tags := r.URL.Query().Get("tags"); tags != "" {
			for _, t := range strings.Split(tags, ",") {
				if t = strings.TrimSpace(t); t != "" {
					filters.Tags = append(filters.Tags, t)
				}
			}
		}
		if ip := r.URL.Query().Get("is_private"); ip == "true" {
			t := true
			filters.IsPrivate = &t
		} else if ip == "false" {
			t := false
			filters.IsPrivate = &t
		}
	}

	capsules, err := CapsuleStore.GetCapsulesByUser(userID, filters)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	page, limit := utils.ParsePagination(r)
	paged, total := utils.SlicePage(capsules, page, limit)
	logger.LogEvent(logger.EventCapsule, r, slog.String("action", "list"), slog.Int("count", len(paged)), slog.Int("total", total))
	utils.JSONPaginatedResponse(w, http.StatusOK, "Capsules fetched", paged, page, limit, total)
}

// CreateCapsule godoc
// @Summary Create capsule
// @Description Create a new capsule
// @Tags capsules
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param input body models.CapsuleInput true "Capsule data"
// @Success 201 {object} models.Capsule
// @Failure 400 {object} map[string]interface{}
// @Router /api/capsules [post]
func CreateCapsule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)
	var req models.CapsuleInput
	if r.Body == nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("empty request body"))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, &utils.ValidationError{Field: "title", Message: "cannot be empty"})
		return
	}
	if len(title) > maxTitleLen {
		utils.ErrorResponse(w, r, http.StatusBadRequest, &utils.ValidationError{Field: "title", Message: "exceeds maximum length"})
		return
	}

	capsule, err := CapsuleStore.AddCapsule(userID, title, req.Content, req.Topic, []string(req.Tags), req.IsPrivate)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}
	logger.LogEvent(logger.EventCapsule, r, slog.String("action", "create"), slog.String("capsule_id", capsule.ID), slog.String("title", title))
	utils.JSONResponse(w, http.StatusCreated, true, "Capsule created", capsule)
}

// CapsuleHandler routes GET/POST to GetCapsules or CreateCapsule.
func CapsuleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetCapsules(w, r)
	case http.MethodPost:
		CreateCapsule(w, r)
	default:
		utils.ErrorResponse(w, r, http.StatusMethodNotAllowed, nil)
	}
}
