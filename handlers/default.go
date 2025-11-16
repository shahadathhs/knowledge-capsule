package handlers

import (
	"net/http"

	"knowledge-capsule-api/utils"
)

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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodGet) {
		return
	}

	utils.JSONResponse(w, http.StatusOK, true, "Service is healthy", map[string]string{
		"service": "knowledge-capsule-api",
	})
}
