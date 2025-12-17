package utils

import (
	"net/http"
	"golang.org/x/time/rate"
)


var globalLimiter = rate.NewLimiter(1000, 50) 

func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		

		if !globalLimiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}