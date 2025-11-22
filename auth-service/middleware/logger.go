package middleware

import (
	"auth-service/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Обрабатываем запрос
		c.Next()

		duration := time.Since(start)

		// Логируем детали запроса
		utils.Log.WithFields(logrus.Fields{
			"client_ip":   c.ClientIP(),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"duration":    duration.String(),
			"user_agent":  c.Request.UserAgent(),
			"timestamp":   time.Now().Format(time.RFC3339),
		}).Info("HTTP request")
	}
}
