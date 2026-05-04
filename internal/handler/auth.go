package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Alexandr20i/workout-tracker/internal/handler/response"
	"github.com/Alexandr20i/workout-tracker/internal/middleware"
	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/Alexandr20i/workout-tracker/internal/repository"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users     *repository.UserRepository
	validate  *validator.Validate
	jwtSecret string
	jwtExpH   int
}

func NewAuthHandler(users *repository.UserRepository, secret string, expH int) *AuthHandler {
	return &AuthHandler{users: users, validate: validator.New(), jwtSecret: secret, jwtExpH: expH}
}

// POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := h.users.FindByEmail(req.Email); err == nil {
		response.Error(w, http.StatusConflict, "email already registered")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := h.users.Create(req.Email, string(hash), req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	token, err := middleware.GenerateToken(user.ID, h.jwtSecret, h.jwtExpH)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response.JSON(w, http.StatusCreated, model.AuthResponse{Token: token, User: *user})
}

// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.users.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		response.Error(w, http.StatusInternalServerError, "server error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := middleware.GenerateToken(user.ID, h.jwtSecret, h.jwtExpH)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response.JSON(w, http.StatusOK, model.AuthResponse{Token: token, User: *user})
}
