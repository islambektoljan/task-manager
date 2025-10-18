package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://redis:6379"
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Failed to parse REDIS_URL: ", err)
	}

	client := redis.NewClient(opts)

	// Test connection with retry
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = client.Ping(ctx).Result()
		cancel()

		if err == nil {
			break
		}

		if i == maxRetries-1 {
			log.Fatal("Failed to connect to Redis after multiple attempts: ", err)
		}

		waitTime := time.Duration(i+1) * 2 * time.Second
		log.Printf("Failed to connect to Redis (attempt %d/%d): %v. Retrying in %v...",
			i+1, maxRetries, err, waitTime)
		time.Sleep(waitTime)
	}

	fmt.Println("Redis connection established")
	return client
}
