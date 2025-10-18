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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Извлечение токена
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		tokenString = strings.TrimSpace(tokenString)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Проверяем, находится ли токен в черном списке
		ctx := context.Background()
		key := fmt.Sprintf("blacklist:%s", tokenString)
		redisExists, err := am.RedisClient.Exists(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if redisExists == 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, userIDExists := claims["user_id"]
		if !userIDExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
			c.Abort()
			return
		}

		// Проверяем expiration time
		exp, expExists := claims["exp"]
		if !expExists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "exp not found in token"})
			c.Abort()
			return
		}

		expFloat, ok := exp.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid exp format in token"})
			c.Abort()
			return
		}

		expTime := int64(expFloat)
		if time.Now().Unix() > expTime {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("token", tokenString)
		c.Next()
	}
}
