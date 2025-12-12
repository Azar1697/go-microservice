package utils

import (
	"net/http"

	"golang.org/x/time/rate"
)

// LimitMiddleware ограничивает количество запросов
func LimitMiddleware(next http.Handler) http.Handler {
	// 1000 запросов в секунду (RPS), "всплеск" (burst) до 5000
	// Burst нужен, чтобы не отбивать запросы, если придет пачка одновременно
	limiter := rate.NewLimiter(1000, 5000)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}