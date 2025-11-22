package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	RedisClient *redis.Client
}

func NewAuthMiddleware(redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{RedisClient: redisClient}
}

func (am *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
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

		// Проверяем, находится ли токен в черном списке
		ctx := context.Background()
		key := fmt.Sprintf("blacklist:%s", tokenString)
		redisExists, err := am.RedisClient.Exists(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error",
				"code":    http.StatusInternalServerError,
			})
			c.Abort()
			return
		}

		if redisExists == 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token has been revoked",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				panic("JWT_SECRET is not set")
			}
			return []byte(secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token",
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

		userID, userIDExists := claims["user_id"]
		if !userIDExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "user_id not found in token",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Извлекаем роль из токена
		role, roleExists := claims["role"]
		if !roleExists {
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

		// Проверяем expiration time
		exp, expExists := claims["exp"]
		if !expExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "exp not found in token",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		expFloat, ok := exp.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid exp format in token",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		expTime := int64(expFloat)
		if time.Now().Unix() > expTime {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token has expired",
				"code":    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Сохраняем данные в контексте
		c.Set("userID", userID)
		c.Set("role", roleStr) // Добавляем роль в контекст
		c.Set("token", tokenString)
		c.Next()
	}
}
