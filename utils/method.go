package utils

import (
	"net/http"
)

func AllowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.Header().Set("Allow", method)
		ErrorResponse(w, http.StatusMethodNotAllowed, nil)
		return false
	}
	return true
}
