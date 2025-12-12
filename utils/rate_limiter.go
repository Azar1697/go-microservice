package utils

import (
	"net/http"

	"golang.org/x/time/rate"
)

// LimitMiddleware ограничивает количество запросов
func LimitMiddleware(next http.Handler) http.Handler {

	limiter := rate.NewLimiter(1000, 5000)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}