package monitoring

import (
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP метрики
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Auth метрики
	AuthAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"type", "status"},
	)

	// Бизнес метрики
	UsersRegistered = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "users_registered_total",
			Help: "Total number of registered users",
		},
	)

	// Системные метрики
	GoroutinesCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_count",
			Help: "Current number of goroutines",
		},
	)
)

// MetricsMiddleware логирует HTTP запросы
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		// Пропускаем метрики endpoint
		if path == "/metrics" {
			c.Next()
			return
		}

		c.Next()

		duration := time.Since(start).Seconds()

		// Обновляем метрики
		HttpRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)

		HttpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		// Обновляем системные метрики
		GoroutinesCount.Set(float64(runtime.NumGoroutine()))
	}
}

// RegisterMetricsHandler регистрирует endpoint для Prometheus
func RegisterMetricsHandler(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// RecordAuthAttempt записывает попытку аутентификации
func RecordAuthAttempt(authType string, success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	AuthAttempts.WithLabelValues(authType, status).Inc()
}

// RecordUserRegistration записывает регистрацию пользователя
func RecordUserRegistration() {
	UsersRegistered.Inc()
}
