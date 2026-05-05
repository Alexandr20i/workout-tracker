package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Alexandr20i/workout-tracker/internal/cache"
	"github.com/Alexandr20i/workout-tracker/internal/handler/response"
	"github.com/Alexandr20i/workout-tracker/internal/middleware"
	"github.com/Alexandr20i/workout-tracker/internal/model"
	"github.com/Alexandr20i/workout-tracker/internal/repository"
	"github.com/redis/go-redis/v9"
)

type StatsHandler struct {
	stats *repository.StatsRepository
	cache *cache.Redis
}

func NewStatsHandler(stats *repository.StatsRepository, cache *cache.Redis) *StatsHandler {
	return &StatsHandler{stats: stats, cache: cache}
}

// Summary godoc
// @Summary      Общая статистика
// @Tags         stats
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.StatsSummary
// @Router       /stats/summary [get]
func (h *StatsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	cacheKey := fmt.Sprintf("stats:summary:%d", userID)

	// Пробуем достать из кэша
	var summary model.StatsSummary
	err := h.cache.Get(r.Context(), cacheKey, &summary)
	if err == nil {
		// Кэш hit — отдаём сразу
		slog.Info("cache hit", "key", cacheKey)
		response.JSON(w, http.StatusOK, summary)
		return
	}

	if !errors.Is(err, redis.Nil) {
		// Неожиданная ошибка Redis — логируем но продолжаем
		slog.Warn("cache error", "error", err)
	}

	// Кэш miss — идём в БД
	slog.Info("cache miss", "key", cacheKey)
	result, err := h.stats.Summary(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch stats")
		return
	}

	// Сохраняем в кэш на 5 минут
	if err := h.cache.Set(context.Background(), cacheKey, result, 5*time.Minute); err != nil {
		slog.Warn("failed to cache stats", "error", err)
	}

	response.JSON(w, http.StatusOK, result)
}

// Progress godoc
// @Summary      Прогресс по упражнению
// @Tags         stats
// @Produce      json
// @Security     BearerAuth
// @Param        exercise_id query int true "ID упражнения"
// @Success      200 {array} model.ProgressPoint
// @Router       /stats/progress [get]
func (h *StatsHandler) Progress(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	exerciseID, err := strconv.ParseInt(r.URL.Query().Get("exercise_id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "exercise_id is required")
		return
	}

	cacheKey := fmt.Sprintf("stats:progress:%d:%d", userID, exerciseID)

	var points []model.ProgressPoint
	err = h.cache.Get(r.Context(), cacheKey, &points)
	if err == nil {
		slog.Info("cache hit", "key", cacheKey)
		response.JSON(w, http.StatusOK, points)
		return
	}

	result, err := h.stats.Progress(userID, exerciseID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch progress")
		return
	}

	if err := h.cache.Set(context.Background(), cacheKey, result, 5*time.Minute); err != nil {
		slog.Warn("failed to cache progress", "error", err)
	}

	response.JSON(w, http.StatusOK, result)
}
