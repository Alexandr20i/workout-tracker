package repository

import "github.com/Alexandr20i/workout-tracker/internal/model"

type UserRepo interface {
	Create(email, passwordHash, name string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByID(id int64) (*model.User, error)
	ListAll() ([]model.User, error)
}

type ExerciseRepo interface {
	Create(userID int64, req *model.CreateExerciseRequest) (*model.Exercise, error)
	ListByUser(userID int64) ([]model.Exercise, error)
	FindByID(id, userID int64) (*model.Exercise, error)
	Delete(id, userID int64) error
	Search(userID int64, query string) ([]model.Exercise, error)
}

type WorkoutRepo interface {
	Create(userID int64, req *model.CreateWorkoutRequest) (*model.Workout, error)
	ListByUser(userID int64) ([]model.Workout, error)
	FindByID(id, userID int64) (*model.Workout, error)
	Delete(id, userID int64) error
}

type SetRepo interface {
	Create(workoutID int64, req *model.CreateSetRequest) (*model.Set, error)
	ListByWorkout(workoutID int64) ([]model.Set, error)
	Delete(id, workoutID int64) error
}

type StatsRepo interface {
	Summary(userID int64) (*model.StatsSummary, error)
	Progress(userID, exerciseID int64) ([]model.ProgressPoint, error)
}
