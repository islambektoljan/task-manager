package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header is required",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Извлечение токена
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		tokenString = strings.TrimSpace(tokenString)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid Authorization header format",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Парсинг JWT токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "your-super-secret-jwt-key-here-make-it-very-long-and-secure"
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token claims",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		userID, exists := claims["user_id"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "user_id not found in token",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid user_id format in token",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Извлекаем роль из токена
		role, exists := claims["role"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "role not found in token",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid role format in token",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		c.Set("userID", userIDStr)
		c.Set("role", roleStr) // Сохраняем роль в контексте
		c.Next()
	}
}
