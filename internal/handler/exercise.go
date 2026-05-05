package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Alexandr20i/workout-tracker/internal/handler/response"
	"github.com/Alexandr20i/workout-tracker/internal/middleware"
	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/Alexandr20i/workout-tracker/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ExerciseHandler struct {
	repo     *repository.ExerciseRepository
	validate *validator.Validate
}

func NewExerciseHandler(repo *repository.ExerciseRepository) *ExerciseHandler {
	return &ExerciseHandler{repo: repo, validate: validator.New()}
}

// GET /exercises
// List godoc
// @Summary      Список упражнений
// @Tags         exercises
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.Exercise
// @Router       /exercises [get]
func (h *ExerciseHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListByUser(middleware.GetUserID(r))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch exercises")
		return
	}
	response.JSON(w, http.StatusOK, list)
}

// POST /exercises
// Create godoc
// @Summary      Создать упражнение
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.CreateExerciseRequest true "Упражнение"
// @Success      201 {object} model.Exercise
// @Router       /exercises [post]
func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	ex, err := h.repo.Create(middleware.GetUserID(r), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create exercise")
		return
	}
	response.JSON(w, http.StatusCreated, ex)
}

// DELETE /exercises/{id}
// Delete godoc
// @Summary      Удалить упражнение
// @Tags         exercises
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID упражнения"
// @Success      200 {object} map[string]string
// @Router       /exercises/{id} [delete]
func (h *ExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(id, middleware.GetUserID(r)); err != nil {
		response.Error(w, http.StatusNotFound, "exercise not found")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
