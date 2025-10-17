package handlers

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"auth-service/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorMsg := "Invalid request data"
		if strings.Contains(err.Error(), "email") {
			errorMsg = "Invalid email format"
		} else if strings.Contains(err.Error(), "password") {
			errorMsg = "Password must be at least 6 characters"
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   errorMsg,
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := h.DB.Where("email = ?", req.Email).First(&existingUser)

	if result.Error == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "User with this email already exists",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Database error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Could not process password",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Create user
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if result := h.DB.Create(&user); result.Error != nil {
		errorMsg := "Could not create user account"
		if strings.Contains(result.Error.Error(), "duplicate key") ||
			strings.Contains(result.Error.Error(), "unique constraint") {
			errorMsg = "User with this email already exists"
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   errorMsg,
				Code:    http.StatusBadRequest,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   errorMsg,
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Could not generate authentication token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := AuthResponse{
		Token:  token,
		UserID: user.ID,
		Email:  user.Email,
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Find user
	var user models.User
	if result := h.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid credentials",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid credentials",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Could not generate token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := AuthResponse{
		Token:  token,
		UserID: user.ID,
		Email:  user.Email,
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    gin.H{"message": "Successfully logged out"},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Convert userID to UUID
	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Generate new token
	token, err := h.generateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Could not generate token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := AuthResponse{
		Token:  token,
		UserID: userID,
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) HealthCheck(c *gin.Context) {
	sqlDB, err := h.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Success: false,
			Error:   "Database connection error",
			Code:    http.StatusServiceUnavailable,
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Success: false,
			Error:   "Database unreachable",
			Code:    http.StatusServiceUnavailable,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    gin.H{"status": "healthy", "service": "auth-service"},
	})
}

func (h *AuthHandler) generateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
