package mocks

import (
	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/stretchr/testify/mock"
)

type ExerciseRepo struct {
	mock.Mock
}

func (m *ExerciseRepo) Create(userID int64, req *model.CreateExerciseRequest) (*model.Exercise, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Exercise), args.Error(1)
}

func (m *ExerciseRepo) ListByUser(userID int64) ([]model.Exercise, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Exercise), args.Error(1)
}

func (m *ExerciseRepo) FindByID(id, userID int64) (*model.Exercise, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Exercise), args.Error(1)
}

func (m *ExerciseRepo) Delete(id, userID int64) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *ExerciseRepo) Search(userID int64, query string) ([]model.Exercise, error) {
	args := m.Called(userID, query)
	return args.Get(0).([]model.Exercise), args.Error(1)
}
