package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alexandr20i/workout-tracker/internal/handler"
	"github.com/Alexandr20i/workout-tracker/internal/middleware"
	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/Alexandr20i/workout-tracker/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper — добавляет userID в контекст запроса (имитирует JWT middleware)
func withUserID(r *http.Request, userID int64) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	return r.WithContext(ctx)
}

func TestExerciseHandler_List(t *testing.T) {
	mockRepo := new(mocks.ExerciseRepo)

	expected := []model.Exercise{
		{ID: 1, UserID: 1, Name: "Планка", MuscleGroup: "Пресс", CreatedAt: time.Now()},
		{ID: 2, UserID: 1, Name: "Приседания", MuscleGroup: "Ноги", CreatedAt: time.Now()},
	}

	// Говорим моку что ожидаем вызов ListByUser(1) и он вернёт expected
	mockRepo.On("ListByUser", int64(1)).Return(expected, nil)

	h := handler.NewExerciseHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/exercises", nil)
	req = withUserID(req, 1)
	rr := httptest.NewRecorder()

	h.List(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string][]model.Exercise
	json.NewDecoder(rr.Body).Decode(&resp)
	assert.Len(t, resp["data"], 2)
	assert.Equal(t, "Планка", resp["data"][0].Name)

	mockRepo.AssertExpectations(t)
}

func TestExerciseHandler_Create_Success(t *testing.T) {
	mockRepo := new(mocks.ExerciseRepo)

	req_body := model.CreateExerciseRequest{
		Name:        "Планка",
		MuscleGroup: "Пресс",
	}
	expected := &model.Exercise{
		ID:          1,
		UserID:      1,
		Name:        "Планка",
		MuscleGroup: "Пресс",
		CreatedAt:   time.Now(),
	}

	mockRepo.On("Create", int64(1), mock.AnythingOfType("*model.CreateExerciseRequest")).
		Return(expected, nil)

	h := handler.NewExerciseHandler(mockRepo)

	body, _ := json.Marshal(req_body)
	req := httptest.NewRequest(http.MethodPost, "/exercises", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserID(req, 1)
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp map[string]*model.Exercise
	json.NewDecoder(rr.Body).Decode(&resp)
	assert.Equal(t, "Планка", resp["data"].Name)

	mockRepo.AssertExpectations(t)
}

func TestExerciseHandler_Create_InvalidJSON(t *testing.T) {
	mockRepo := new(mocks.ExerciseRepo)
	h := handler.NewExerciseHandler(mockRepo)

	req := httptest.NewRequest(http.MethodPost, "/exercises", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	req = withUserID(req, 1)
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	// Мок не должен был вызываться
	mockRepo.AssertNotCalled(t, "Create")
}

func TestExerciseHandler_Create_ValidationFail(t *testing.T) {
	mockRepo := new(mocks.ExerciseRepo)
	h := handler.NewExerciseHandler(mockRepo)

	// Пустое имя — должна упасть валидация
	body, _ := json.Marshal(map[string]string{"name": ""})
	req := httptest.NewRequest(http.MethodPost, "/exercises", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withUserID(req, 1)
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockRepo.AssertNotCalled(t, "Create")
}
