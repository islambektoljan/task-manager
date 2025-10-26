package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет наличие X-User-ID заголовка (переданного через API Gateway)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "X-User-ID header is required",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
