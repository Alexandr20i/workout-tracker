package mocks

import (
	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/stretchr/testify/mock"
)

type SetRepo struct {
	mock.Mock
}

func (m *SetRepo) Create(workoutID int64, req *model.CreateSetRequest) (*model.Set, error) {
	args := m.Called(workoutID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Set), args.Error(1)
}

func (m *SetRepo) ListByWorkout(workoutID int64) ([]model.Set, error) {
	args := m.Called(workoutID)
	return args.Get(0).([]model.Set), args.Error(1)
}

func (m *SetRepo) Delete(id, workoutID int64) error {
	args := m.Called(id, workoutID)
	return args.Error(0)
}
