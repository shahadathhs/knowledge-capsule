package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

func ParseAndValidateBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	// Decode JSON body
	if r.Body == nil {
		ErrorResponse(w, r, http.StatusBadRequest, errors.New("empty request body"))
		return false
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		ErrorResponse(w, r, http.StatusBadRequest, err)
		return false
	}

	// Basic validation: check for empty string fields
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// If the field is a string and empty
		if field.Kind() == reflect.String && field.String() == "" {
			ErrorResponse(w, r, http.StatusBadRequest,
				&ValidationError{Field: fieldType.Name, Message: "cannot be empty"})
			return false
		}
	}

	return true
}

// ValidationError represents a simple validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *ValidationError) Error() string {
	return v.Field + ": " + v.Message
}
