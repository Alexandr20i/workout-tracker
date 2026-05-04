package model

import "time"

type User struct {
	ID           int64     `db:"id"            json:"id"`
	Email        string    `db:"email"         json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name"          json:"name"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}

type Exercise struct {
	ID          int64     `db:"id"           json:"id"`
	UserID      int64     `db:"user_id"      json:"user_id"`
	Name        string    `db:"name"         json:"name"`
	Description string    `db:"description"  json:"description"`
	MuscleGroup string    `db:"muscle_group" json:"muscle_group"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}

type Workout struct {
	ID        int64     `db:"id"         json:"id"`
	UserID    int64     `db:"user_id"    json:"user_id"`
	Title     string    `db:"title"      json:"title"`
	Notes     string    `db:"notes"      json:"notes"`
	Date      time.Time `db:"date"       json:"date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	Sets []Set `db:"-" json:"sets,omitempty"`
}

type Set struct {
	ID           int64     `db:"id"            json:"id"`
	WorkoutID    int64     `db:"workout_id"    json:"workout_id"`
	ExerciseID   int64     `db:"exercise_id"   json:"exercise_id"`
	Reps         int       `db:"reps"          json:"reps"`
	WeightKg     float64   `db:"weight_kg"     json:"weight_kg"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	ExerciseName string    `db:"exercise_name" json:"exercise_name,omitempty"`
}

// --- Request DTOs ---

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name"     validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateExerciseRequest struct {
	Name        string `json:"name"         validate:"required"`
	Description string `json:"description"`
	MuscleGroup string `json:"muscle_group"`
}

type CreateWorkoutRequest struct {
	Title string `json:"title" validate:"required"`
	Notes string `json:"notes"`
	Date  string `json:"date"` // YYYY-MM-DD
}

type CreateSetRequest struct {
	ExerciseID int64   `json:"exercise_id" validate:"required"`
	Reps       int     `json:"reps"        validate:"required,min=1"`
	WeightKg   float64 `json:"weight_kg"   validate:"min=0"`
}

// --- Response DTOs ---

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ProgressPoint struct {
	Date      time.Time `db:"date"       json:"date"`
	MaxWeight float64   `db:"max_weight" json:"max_weight"`
	TotalReps int       `db:"total_reps" json:"total_reps"`
}

type StatsSummary struct {
	TotalWorkouts    int     `db:"total_workouts"     json:"total_workouts"`
	TotalSets        int     `db:"total_sets"         json:"total_sets"`
	TotalWeightKg    float64 `db:"total_weight_kg"    json:"total_weight_kg"`
	WorkoutsThisWeek int     `db:"workouts_this_week" json:"workouts_this_week"`
}
