package utils

import (
	"encoding/json"
	"net/http"

	"knowledge-capsule-api/models"
)

func JSONResponse(w http.ResponseWriter, status int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := models.APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, status int, err error) {
	JSONResponse(w, status, false, "", map[string]string{"error": err.Error()})
}
