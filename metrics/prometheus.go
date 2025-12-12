package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Счетчик запросов (RPS)
var TotalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests processed, labeled by path.",
	},
	[]string{"path"},
)

// Гистограмма времени ответа (Latency) - ТРЕБОВАНИЕ КРИТЕРИЯ 2
var HttpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.DefBuckets, // Стандартные бакеты времени (0.005, 0.01, 0.025 ...)
	},
	[]string{"path"},
)

func init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(HttpDuration) // Регистрируем новую метрику
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Засекаем время
		
		next.ServeHTTP(w, r)
		
		duration := time.Since(start).Seconds() // Считаем длительность
		
		// Записываем данные
		TotalRequests.WithLabelValues(r.URL.Path).Inc()
		HttpDuration.WithLabelValues(r.URL.Path).Observe(duration)
	})
}