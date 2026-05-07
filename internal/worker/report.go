package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/Alexandr20i/workout-tracker/internal/repository"
)

type ReportWorker struct {
	statsRepo *repository.StatsRepository
	userRepo  *repository.UserRepository
	interval  time.Duration
	stopCh    chan struct{}
}

func NewReportWorker(
	statsRepo *repository.StatsRepository,
	userRepo *repository.UserRepository,
	interval time.Duration,
) *ReportWorker {
	return &ReportWorker{
		statsRepo: statsRepo,
		userRepo:  userRepo,
		interval:  interval,
		stopCh:    make(chan struct{}),
	}
}

// Start запускает воркер в отдельной горутине
func (w *ReportWorker) Start(ctx context.Context) {
	go w.run(ctx)
	slog.Info("report worker started", "interval", w.interval)
}

// Stop останавливает воркер
func (w *ReportWorker) Stop() {
	close(w.stopCh)
	slog.Info("report worker stopped")
}

func (w *ReportWorker) run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.generateReports(ctx)
		case <-w.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (w *ReportWorker) generateReports(ctx context.Context) {
	slog.Info("generating weekly reports...")

	// Получаем всех пользователей
	users, err := w.userRepo.ListAll()
	if err != nil {
		slog.Error("failed to fetch users for report", "error", err)
		return
	}

	for _, user := range users {
		summary, err := w.statsRepo.Summary(user.ID)
		if err != nil {
			slog.Error("failed to get summary", "user_id", user.ID, "error", err)
			continue
		}

		// В реальном проекте — отправляем email
		// Здесь логируем как демонстрацию паттерна
		slog.Info("weekly report",
			"user_id", user.ID,
			"user_name", user.Name,
			"total_workouts", summary.TotalWorkouts,
			"workouts_this_week", summary.WorkoutsThisWeek,
			"total_weight_kg", summary.TotalWeightKg,
		)
	}

	slog.Info("weekly reports done", "users_processed", len(users))
}
