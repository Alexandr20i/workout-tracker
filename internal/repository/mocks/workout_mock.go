package mocks

import (
	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/stretchr/testify/mock"
)

type WorkoutRepo struct {
	mock.Mock
}

func (m *WorkoutRepo) Create(userID int64, req *model.CreateWorkoutRequest) (*model.Workout, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Workout), args.Error(1)
}

func (m *WorkoutRepo) ListByUser(userID int64) ([]model.Workout, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Workout), args.Error(1)
}

func (m *WorkoutRepo) FindByID(id, userID int64) (*model.Workout, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Workout), args.Error(1)
}

func (m *WorkoutRepo) Delete(id, userID int64) error {
	args := m.Called(id, userID)
	return args.Error(0)
}
