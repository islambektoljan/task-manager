package monitoring

import (
	"expvar"
	"net/http"
	"runtime"
	"time"
)

var (
	startTime = time.Now().UTC()

	// Кастомные expvar метрики
	Requests = expvar.NewInt("http_requests")
	Errors   = expvar.NewInt("http_errors")
)

func InitExpvar() {
	// Стандартные метрики Go
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("uptime", expvar.Func(func() interface{} {
		return time.Since(startTime).Seconds()
	}))

	// Кастомные метрики приложения
	expvar.Publish("version", expvar.Func(func() interface{} {
		return "1.0.0"
	}))

	// Инициализация счетчиков
	Requests.Set(0)
	Errors.Set(0)
}

func ExpvarHandler() http.Handler {
	return expvar.Handler()
}

func RecordRequest() {
	Requests.Add(1)
}

func RecordError() {
	Errors.Add(1)
}
