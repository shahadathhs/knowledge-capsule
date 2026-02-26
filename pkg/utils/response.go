package utils

import (
	"encoding/json"
	"net/http"

	"knowledge-capsule/app/models"
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

// PaginatedResponse is the structure for paginated list responses.
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int         `json:"total"`
}

// JSONPaginatedResponse writes a paginated JSON response.
func JSONPaginatedResponse(w http.ResponseWriter, status int, message string, data interface{}, page, limit, total int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Page:    page,
		Limit:   limit,
		Total:   total,
	}

	json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, status int, err error) {
	var errorMessage string

	switch {
	case err != nil:
		errorMessage = err.Error()
	case status >= 500:
		errorMessage = "Internal server error"
	case status == http.StatusUnauthorized:
		errorMessage = "Unauthorized"
	case status == http.StatusBadRequest:
		errorMessage = "Bad request"
	default:
		errorMessage = "An error occurred"
	}

	JSONResponse(w, status, false, "", map[string]string{"error": errorMessage})
}
