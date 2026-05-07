package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alexandr20i/workout-tracker/internal/handler"
	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/Alexandr20i/workout-tracker/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkoutHandler_Create_Success(t *testing.T) {
	mockWorkouts := new(mocks.WorkoutRepo)
	mockSets := new(mocks.SetRepo)

	expected := &model.Workout{
		ID:        1,
		UserID:    1,
		Title:     "Тренировка ног",
		Date:      time.Now(),
		CreatedAt: time.Now(),
	}

	mockWorkouts.On("Create", int64(1), mock.AnythingOfType("*model.CreateWorkoutRequest")).
		Return(expected, nil)

	h := handler.NewWorkoutHandler(mockWorkouts, mockSets)

	body, _ := json.Marshal(model.CreateWorkoutRequest{Title: "Тренировка ног"})
	req := httptest.NewRequest(http.MethodPost, "/workouts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserID(req, 1)
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockWorkouts.AssertExpectations(t)
}

func TestWorkoutHandler_List(t *testing.T) {
	mockWorkouts := new(mocks.WorkoutRepo)
	mockSets := new(mocks.SetRepo)

	expected := []model.Workout{
		{ID: 1, UserID: 1, Title: "Тренировка ног", CreatedAt: time.Now()},
	}
	mockWorkouts.On("ListByUser", int64(1)).Return(expected, nil)

	h := handler.NewWorkoutHandler(mockWorkouts, mockSets)

	req := httptest.NewRequest(http.MethodGet, "/workouts", nil)
	req = withUserID(req, 1)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string][]model.Workout
	json.NewDecoder(rr.Body).Decode(&resp)
	assert.Len(t, resp["data"], 1)

	mockWorkouts.AssertExpectations(t)
}

func TestWorkoutHandler_Create_MissingTitle(t *testing.T) {
	mockWorkouts := new(mocks.WorkoutRepo)
	mockSets := new(mocks.SetRepo)

	h := handler.NewWorkoutHandler(mockWorkouts, mockSets)

	body, _ := json.Marshal(map[string]string{"title": ""})
	req := httptest.NewRequest(http.MethodPost, "/workouts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserID(req, 1)
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockWorkouts.AssertNotCalled(t, "Create")
}
