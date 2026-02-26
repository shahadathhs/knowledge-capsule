package handlers

import (
	"net/http"

	"knowledge-capsule/pkg/utils"
)

// TestChatHandler serves the WebSocket chat test page at /test-ws.
func TestChatHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}
	http.ServeFile(w, r, "web/test_chat.html")
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, "Welcome to Knowledge Capsule API", map[string]string{
		"status": "ok",
	})
}

func ApiRootHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, "Knowledge Capsule API Root", map[string]string{
		"version": "v1",
	})
}

// HealthHandler godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, "Service is healthy", map[string]string{
		"service": "knowledge-capsule",
	})
}
