package middleware

import (
	"net/http"
	"sync"

	"github.com/Alexandr20i/workout-tracker/internal/handler/response"
	"golang.org/x/time/rate"
)

// ipLimiter хранит лимитер для каждого IP
type ipLimiter struct {
	limiter *rate.Limiter
}

var (
	limiters = make(map[string]*ipLimiter)
	mu       sync.Mutex
)

// getLimiter возвращает лимитер для IP, создаёт если нет
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if l, ok := limiters[ip]; ok {
		return l.limiter
	}

	// 10 запросов в секунду, burst до 20
	l := rate.NewLimiter(10, 20)
	limiters[ip] = &ipLimiter{limiter: l}
	return l
}

// RateLimit — middleware ограничивает количество запросов с одного IP
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := getLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			response.Error(w, http.StatusTooManyRequests, "too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}
