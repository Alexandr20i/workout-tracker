package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexandr20i/workout-tracker/config"
	"github.com/alexandr20i/workout-tracker/internal/handler"
	"github.com/alexandr20i/workout-tracker/internal/middleware"
	"github.com/alexandr20i/workout-tracker/internal/repository"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := sqlx.Connect("postgres", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer db.Close()
	log.Println("✅ Connected to PostgreSQL")

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	exerciseRepo := repository.NewExerciseRepository(db)
	workoutRepo := repository.NewWorkoutRepository(db)
	setRepo := repository.NewSetRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	// Хендлеры
	authHandler := handler.NewAuthHandler(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	exerciseHandler := handler.NewExerciseHandler(exerciseRepo)
	workoutHandler := handler.NewWorkoutHandler(workoutRepo, setRepo)
	statsHandler := handler.NewStatsHandler(statsRepo)

	// Роутер
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Публичные маршруты
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	// Защищённые маршруты
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWT.Secret))

		r.Get("/exercises", exerciseHandler.List)
		r.Post("/exercises", exerciseHandler.Create)
		r.Delete("/exercises/{id}", exerciseHandler.Delete)

		r.Get("/workouts", workoutHandler.List)
		r.Post("/workouts", workoutHandler.Create)
		r.Get("/workouts/{id}", workoutHandler.Get)
		r.Delete("/workouts/{id}", workoutHandler.Delete)

		r.Post("/workouts/{id}/sets", workoutHandler.AddSet)
		r.Delete("/workouts/{id}/sets/{setId}", workoutHandler.DeleteSet)

		r.Get("/stats/summary", statsHandler.Summary)
		r.Get("/stats/progress", statsHandler.Progress)
	})

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 Server running on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
