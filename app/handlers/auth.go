package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"knowledge-capsule/pkg/logger"
	"knowledge-capsule/pkg/utils"
)

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with name, email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body object{name=string,email=string,password=string} true "User registration info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodPost) {
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if !utils.ParseAndValidateBody(w, r, &req) {
		return
	}

	user, err := UserStore.AddUser(req.Name, req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, err)
		return
	}
	logger.LogEvent(logger.EventAuth, r, slog.String("action", "register"), slog.String("user_id", user.ID), slog.String("email", req.Email))

	utils.JSONResponse(w, http.StatusCreated, true, "User registered", map[string]string{
		"user_id": user.ID,
	})
}

// LoginHandler godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body object{email=string,password=string} true "User login info"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.AllowMethod(w, r, http.MethodPost) {
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !utils.ParseAndValidateBody(w, r, &req) {
		return
	}

	user, err := UserStore.FindByEmail(req.Email)
	if err != nil || user == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	role := user.Role
	if role == "" {
		role = "user"
	}
	token, err := utils.GenerateJWT(user.ID, user.Email, role, time.Hour*24)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	logger.LogEvent(logger.EventAuth, r, slog.String("action", "login"), slog.String("user_id", user.ID), slog.String("email", req.Email))

	utils.JSONResponse(w, http.StatusOK, true, "Login successful", map[string]string{
		"token": token,
	})
}
