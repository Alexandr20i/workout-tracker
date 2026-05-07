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

type WorkoutHandler struct {
	workouts repository.WorkoutRepo
	sets     repository.SetRepo
	validate *validator.Validate
}

func NewWorkoutHandler(w *repository.WorkoutRepository, s *repository.SetRepository) *WorkoutHandler {
	return &WorkoutHandler{workouts: w, sets: s, validate: validator.New()}
}

// GET /workouts
// List godoc
// @Summary      История тренировок
// @Tags         workouts
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} model.Workout
// @Router       /workouts [get]
func (h *WorkoutHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.workouts.ListByUser(middleware.GetUserID(r))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch workouts")
		return
	}
	response.JSON(w, http.StatusOK, list)
}

// POST /workouts
// Create godoc
// @Summary      Создать тренировку
// @Tags         workouts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body model.CreateWorkoutRequest true "Тренировка"
// @Success      201 {object} model.Workout
// @Router       /workouts [post]
func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	workout, err := h.workouts.Create(middleware.GetUserID(r), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, workout)
}

// GET /workouts/{id}
// Get godoc
// @Summary      Детали тренировки
// @Tags         workouts
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID тренировки"
// @Success      200 {object} model.Workout
// @Router       /workouts/{id} [get]
func (h *WorkoutHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	workout, err := h.workouts.FindByID(id, userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "workout not found")
		return
	}

	sets, err := h.sets.ListByWorkout(id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch sets")
		return
	}
	workout.Sets = sets
	response.JSON(w, http.StatusOK, workout)
}

// DELETE /workouts/{id}
// Delete godoc
// @Summary      Удалить тренировку
// @Tags         workouts
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID тренировки"
// @Success      200 {object} map[string]string
// @Router       /workouts/{id} [delete]

func (h *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.workouts.Delete(id, middleware.GetUserID(r)); err != nil {
		response.Error(w, http.StatusNotFound, "workout not found")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// POST /workouts/{id}/sets
// AddSet godoc
// @Summary      Добавить подход
// @Tags         workouts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path int true "ID тренировки"
// @Param        body body model.CreateSetRequest true "Подход"
// @Success      201 {object} model.Set
// @Router       /workouts/{id}/sets [post]
func (h *WorkoutHandler) AddSet(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	workoutID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid workout id")
		return
	}

	// Проверяем владельца
	if _, err := h.workouts.FindByID(workoutID, userID); err != nil {
		response.Error(w, http.StatusNotFound, "workout not found")
		return
	}

	var req model.CreateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	set, err := h.sets.Create(workoutID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create set")
		return
	}
	response.JSON(w, http.StatusCreated, set)
}

// DELETE /workouts/{id}/sets/{setId}
// DeleteSet godoc
// @Summary      Удалить подход
// @Tags         workouts
// @Produce      json
// @Security     BearerAuth
// @Param        id    path int true "ID тренировки"
// @Param        setId path int true "ID подхода"
// @Success      200 {object} map[string]string
// @Router       /workouts/{id}/sets/{setId} [delete]
func (h *WorkoutHandler) DeleteSet(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	workoutID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	setID, err := strconv.ParseInt(chi.URLParam(r, "setId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid set id")
		return
	}
	if _, err := h.workouts.FindByID(workoutID, userID); err != nil {
		response.Error(w, http.StatusNotFound, "workout not found")
		return
	}
	if err := h.sets.Delete(setID, workoutID); err != nil {
		response.Error(w, http.StatusNotFound, "set not found")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
