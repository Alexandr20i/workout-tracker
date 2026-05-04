package handler

import (
	"net/http"
	"strconv"

	"github.com/Alexandr20i/workout-tracker/internal/handler/response"
	"github.com/Alexandr20i/workout-tracker/internal/middleware"
	"github.com/Alexandr20i/workout-tracker/internal/repository"
)

type StatsHandler struct {
	stats *repository.StatsRepository
}

func NewStatsHandler(stats *repository.StatsRepository) *StatsHandler {
	return &StatsHandler{stats: stats}
}

// GET /stats/summary
func (h *StatsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.stats.Summary(middleware.GetUserID(r))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch stats")
		return
	}
	response.JSON(w, http.StatusOK, summary)
}

// GET /stats/progress?exercise_id=1
func (h *StatsHandler) Progress(w http.ResponseWriter, r *http.Request) {
	exerciseID, err := strconv.ParseInt(r.URL.Query().Get("exercise_id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "exercise_id is required")
		return
	}
	points, err := h.stats.Progress(middleware.GetUserID(r), exerciseID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch progress")
		return
	}
	response.JSON(w, http.StatusOK, points)
}
