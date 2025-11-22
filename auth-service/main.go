package main

import (
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middleware"
	"auth-service/monitoring"
	"auth-service/utils"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		utils.Log.Warn("No .env file found")
	}

	// Initialize logger
	utils.InitLogger()
	utils.Log.Info("Starting auth-service...")

	// Initialize expvar metrics
	monitoring.InitExpvar()

	// Connect to database
	db := database.ConnectDB()

	// Connect to Redis
	redisClient := database.ConnectRedis()
	defer func() {
		if err := redisClient.Close(); err != nil {
			utils.Log.Errorf("Error closing Redis: %v", err)
		}
	}()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		utils.Log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create router
	r := gin.Default()

	// CORS configuration - РАСКОММЕНТИРУЙТЕ ЭТО!
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create auth handler and middleware
	authHandler := handlers.NewAuthHandler(db, redisClient)
	authMiddleware := middleware.NewAuthMiddleware(redisClient)

	// Security and monitoring middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RequestLogger())
	r.Use(monitoring.MetricsMiddleware())

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// Metrics endpoints
	monitoring.RegisterMetricsHandler(r)
	r.GET("/debug/vars", gin.WrapH(monitoring.ExpvarHandler())) // Добавьте этот endpoint

	// Health check endpoint
	r.GET("/health", authHandler.HealthCheck)

	// Protected routes
	auth := r.Group("/")
	auth.Use(authMiddleware.AuthMiddleware())
	{
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Create server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		utils.Log.Infof("Auth service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.Log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log.Info("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Log.Fatalf("Server forced to shutdown: %v", err)
	}

	utils.Log.Info("Server exited")
}
