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

// GetCapsuleByID godoc
// @Summary Get capsule by ID
// @Description Get a single capsule (user must own it)
// @Tags capsules
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Capsule ID"
// @Success 200 {object} models.Capsule
// @Failure 404 {object} map[string]interface{}
// @Router /api/capsules/{id} [get]
func GetCapsuleByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)
	id := strings.TrimPrefix(r.URL.Path, "/api/capsules/")
	if id == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("missing capsule id"))
		return
	}
	capsule, err := CapsuleStore.FindByID(id)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusNotFound, err)
		return
	}
	if capsule.UserID != userID {
		utils.ErrorResponse(w, r, http.StatusNotFound, errors.New("capsule not found"))
		return
	}
	utils.JSONResponse(w, http.StatusOK, true, "Capsule fetched", capsule)
}

// UpdateCapsule godoc
// @Summary Update capsule by ID
// @Description Update a capsule (user must own it)
// @Tags capsules
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Capsule ID"
// @Param input body models.CapsuleInput true "Updated capsule fields"
// @Success 200 {object} models.Capsule
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/capsules/{id} [put]
func UpdateCapsule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)
	id := strings.TrimPrefix(r.URL.Path, "/api/capsules/")
	if id == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("missing capsule id"))
		return
	}
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

	updated := models.Capsule{CapsuleInput: req}
	capsule, err := CapsuleStore.UpdateCapsule(id, userID, updated)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusNotFound, err)
		return
	}
	logger.LogEvent(logger.EventCapsule, r, slog.String("action", "update"), slog.String("capsule_id", id))
	utils.JSONResponse(w, http.StatusOK, true, "Capsule updated", capsule)
}

// DeleteCapsule godoc
// @Summary Delete capsule by ID
// @Description Delete a capsule (user must own it)
// @Tags capsules
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Capsule ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/capsules/{id} [delete]
func DeleteCapsule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserContextKey).(string)
	id := strings.TrimPrefix(r.URL.Path, "/api/capsules/")
	if id == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, errors.New("missing capsule id"))
		return
	}
	if err := CapsuleStore.DeleteCapsule(id, userID); err != nil {
		utils.ErrorResponse(w, r, http.StatusNotFound, err)
		return
	}
	logger.LogEvent(logger.EventCapsule, r, slog.String("action", "delete"), slog.String("capsule_id", id))
	utils.JSONResponse(w, http.StatusOK, true, "Capsule deleted", nil)
}

// CapsuleByIDHandler routes GET/PUT/DELETE to the appropriate handler.
func CapsuleByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetCapsuleByID(w, r)
	case http.MethodPut:
		UpdateCapsule(w, r)
	case http.MethodDelete:
		DeleteCapsule(w, r)
	default:
		utils.ErrorResponse(w, r, http.StatusMethodNotAllowed, nil)
	}
}
