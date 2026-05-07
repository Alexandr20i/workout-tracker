package repository

import (
	"fmt"
	"time"

	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/jmoiron/sqlx"
)

// ===================== User =====================

type UserRepository struct{ db *sqlx.DB }

func NewUserRepository(db *sqlx.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(email, passwordHash, name string) (*model.User, error) {
	u := &model.User{}
	return u, r.db.QueryRowx(`
		INSERT INTO users (email, password_hash, name)
		VALUES ($1, $2, $3)
		RETURNING *`,
		email, passwordHash, name,
	).StructScan(u)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	return u, r.db.QueryRowx(`SELECT * FROM users WHERE email = $1`, email).StructScan(u)
}

func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	u := &model.User{}
	return u, r.db.QueryRowx(`SELECT * FROM users WHERE id = $1`, id).StructScan(u)
}

func (r *UserRepository) ListAll() ([]model.User, error) {
	var users []model.User
	return users, r.db.Select(&users, `SELECT * FROM users ORDER BY id`)
}

// ===================== Exercise =====================

type ExerciseRepository struct{ db *sqlx.DB }

func NewExerciseRepository(db *sqlx.DB) *ExerciseRepository { return &ExerciseRepository{db: db} }

func (r *ExerciseRepository) Create(userID int64, req *model.CreateExerciseRequest) (*model.Exercise, error) {
	ex := &model.Exercise{}
	return ex, r.db.QueryRowx(`
		INSERT INTO exercises (user_id, name, description, muscle_group)
		VALUES ($1, $2, $3, $4) RETURNING *`,
		userID, req.Name, req.Description, req.MuscleGroup,
	).StructScan(ex)
}

func (r *ExerciseRepository) ListByUser(userID int64) ([]model.Exercise, error) {
	var list []model.Exercise
	return list, r.db.Select(&list, `SELECT * FROM exercises WHERE user_id = $1 ORDER BY name`, userID)
}

func (r *ExerciseRepository) FindByID(id, userID int64) (*model.Exercise, error) {
	ex := &model.Exercise{}
	return ex, r.db.QueryRowx(`SELECT * FROM exercises WHERE id = $1 AND user_id = $2`, id, userID).StructScan(ex)
}

func (r *ExerciseRepository) Delete(id, userID int64) error {
	res, err := r.db.Exec(`DELETE FROM exercises WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

// Search — полнотекстовый поиск по упражнениям
func (r *ExerciseRepository) Search(userID int64, query string) ([]model.Exercise, error) {
	var list []model.Exercise
	return list, r.db.Select(&list, `
		SELECT *
		FROM exercises
		WHERE user_id = $1
		  AND search_vector @@ plainto_tsquery('russian', $2)
		ORDER BY ts_rank(search_vector, plainto_tsquery('russian', $2)) DESC`,
		userID, query,
	)
}

// ===================== Workout =====================

type WorkoutRepository struct{ db *sqlx.DB }

func NewWorkoutRepository(db *sqlx.DB) *WorkoutRepository { return &WorkoutRepository{db: db} }

func (r *WorkoutRepository) Create(userID int64, req *model.CreateWorkoutRequest) (*model.Workout, error) {
	w := &model.Workout{}
	date := time.Now()
	if req.Date != "" {
		var err error
		date, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
		}
	}
	return w, r.db.QueryRowx(`
		INSERT INTO workouts (user_id, title, notes, date)
		VALUES ($1, $2, $3, $4) RETURNING *`,
		userID, req.Title, req.Notes, date,
	).StructScan(w)
}

func (r *WorkoutRepository) ListByUser(userID int64) ([]model.Workout, error) {
	var list []model.Workout
	return list, r.db.Select(&list, `SELECT * FROM workouts WHERE user_id = $1 ORDER BY date DESC`, userID)
}

func (r *WorkoutRepository) FindByID(id, userID int64) (*model.Workout, error) {
	w := &model.Workout{}
	return w, r.db.QueryRowx(`SELECT * FROM workouts WHERE id = $1 AND user_id = $2`, id, userID).StructScan(w)
}

func (r *WorkoutRepository) Delete(id, userID int64) error {
	res, err := r.db.Exec(`DELETE FROM workouts WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

// ===================== Set =====================

type SetRepository struct{ db *sqlx.DB }

func NewSetRepository(db *sqlx.DB) *SetRepository { return &SetRepository{db: db} }

func (r *SetRepository) Create(workoutID int64, req *model.CreateSetRequest) (*model.Set, error) {
	s := &model.Set{}
	return s, r.db.QueryRowx(`
		INSERT INTO sets (workout_id, exercise_id, reps, weight_kg)
		VALUES ($1, $2, $3, $4) RETURNING *`,
		workoutID, req.ExerciseID, req.Reps, req.WeightKg,
	).StructScan(s)
}

func (r *SetRepository) ListByWorkout(workoutID int64) ([]model.Set, error) {
	var list []model.Set
	return list, r.db.Select(&list, `
		SELECT s.*, e.name AS exercise_name
		FROM sets s
		JOIN exercises e ON e.id = s.exercise_id
		WHERE s.workout_id = $1
		ORDER BY s.created_at`,
		workoutID,
	)
}

func (r *SetRepository) Delete(id, workoutID int64) error {
	res, err := r.db.Exec(`DELETE FROM sets WHERE id = $1 AND workout_id = $2`, id, workoutID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

// ===================== Stats =====================

type StatsRepository struct{ db *sqlx.DB }

func NewStatsRepository(db *sqlx.DB) *StatsRepository { return &StatsRepository{db: db} }

func (r *StatsRepository) Summary(userID int64) (*model.StatsSummary, error) {
	s := &model.StatsSummary{}
	return s, r.db.QueryRowx(`
		SELECT
			COUNT(DISTINCT w.id)                    AS total_workouts,
			COUNT(s.id)                             AS total_sets,
			COALESCE(SUM(s.weight_kg * s.reps), 0) AS total_weight_kg,
			COUNT(DISTINCT w.id) FILTER (
				WHERE w.date >= date_trunc('week', CURRENT_DATE)
			)                                       AS workouts_this_week
		FROM workouts w
		LEFT JOIN sets s ON s.workout_id = w.id
		WHERE w.user_id = $1`,
		userID,
	).StructScan(s)
}

func (r *StatsRepository) Progress(userID, exerciseID int64) ([]model.ProgressPoint, error) {
	var list []model.ProgressPoint
	return list, r.db.Select(&list, `
		SELECT
			w.date,
			MAX(s.weight_kg) AS max_weight,
			SUM(s.reps)      AS total_reps
		FROM sets s
		JOIN workouts w ON w.id = s.workout_id
		WHERE w.user_id = $1 AND s.exercise_id = $2
		GROUP BY w.date
		ORDER BY w.date`,
		userID, exerciseID,
	)
}
