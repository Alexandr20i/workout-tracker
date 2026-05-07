package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Alexandr20i/workout-tracker/config"
	_ "github.com/Alexandr20i/workout-tracker/docs"
	"github.com/Alexandr20i/workout-tracker/internal/cache"
	"github.com/Alexandr20i/workout-tracker/internal/handler"
	"github.com/Alexandr20i/workout-tracker/internal/middleware"
	"github.com/Alexandr20i/workout-tracker/internal/repository"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           Workout Tracker API
// @version         1.0
// @description     REST API для отслеживания тренировок
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := sqlx.Connect("postgres", cfg.DB.DSN())
	if err != nil {
		slog.Error("db connect error", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("connected to PostgreSQL")

	// Redis
	redisClient, err := cache.NewRedis(cfg.Redis.Addr)
	if err != nil {
		slog.Error("redis connect error", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to Redis")

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
	statsHandler := handler.NewStatsHandler(statsRepo, redisClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWT.Secret))

		r.Get("/exercises", exerciseHandler.List)
		r.Get("/exercises/search", exerciseHandler.Search) // <- новый
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
	slog.Info("server started", "addr", "http://localhost"+addr)
	slog.Info("swagger UI", "addr", "http://localhost"+addr+"/swagger/")
	log.Fatal(http.ListenAndServe(addr, r))
}
