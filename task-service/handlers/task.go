package handlers

import (
	"net/http"
	"time"

	"task-service/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

func (h *TaskHandler) HealthCheck(c *gin.Context) {
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
		Data:    gin.H{"status": "healthy", "service": "task-service"},
	})
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	role, _ := c.Get("role") // Получаем роль

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var tasks []models.Task

	// Если пользователь администратор, показываем все задачи, иначе только свои
	if role == "admin" {
		result := h.DB.Find(&tasks)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   "Failed to fetch tasks",
				Code:    http.StatusInternalServerError,
			})
			return
		}
	} else {
		result := h.DB.Where("created_by = ?", userUUID).Find(&tasks)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   "Failed to fetch tasks",
				Code:    http.StatusInternalServerError,
			})
			return
		}
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    tasks,
	})
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	role, _ := c.Get("role")

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	taskID := c.Param("id")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid task ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var task models.Task
	var result *gorm.DB

	if role == "admin" {
		result = h.DB.Where("id = ?", taskUUID).First(&task)
	} else {
		result = h.DB.Where("id = ? AND created_by = ?", taskUUID, userUUID).First(&task)
	}

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error:   "Task not found",
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to fetch task",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    task,
	})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Валидация статуса
	if req.Status != "" {
		validStatus := map[string]bool{
			string(models.StatusPending):    true,
			string(models.StatusInProgress): true,
			string(models.StatusCompleted):  true,
			string(models.StatusCancelled):  true,
		}
		if !validStatus[req.Status] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Invalid status value",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	// Валидация приоритета
	if req.Priority != "" {
		validPriority := map[string]bool{
			string(models.PriorityLow):    true,
			string(models.PriorityMedium): true,
			string(models.PriorityHigh):   true,
			string(models.PriorityUrgent): true,
		}
		if !validPriority[req.Priority] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Invalid priority value",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      models.TaskStatus(req.Status),
		Priority:    models.TaskPriority(req.Priority),
		DueDate:     req.DueDate,
		CreatedBy:   userUUID,
	}

	result := h.DB.Create(&task)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to create task",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    task,
	})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	role, _ := c.Get("role")

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	taskID := c.Param("id")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid task ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Находим задачу
	var task models.Task
	var result *gorm.DB

	if role == "admin" {
		result = h.DB.Where("id = ?", taskUUID).First(&task)
	} else {
		result = h.DB.Where("id = ? AND created_by = ?", taskUUID, userUUID).First(&task)
	}

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Error:   "Task not found",
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to fetch task",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Обновляем поля
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		// Валидация статуса
		validStatus := map[string]bool{
			string(models.StatusPending):    true,
			string(models.StatusInProgress): true,
			string(models.StatusCompleted):  true,
			string(models.StatusCancelled):  true,
		}
		if !validStatus[req.Status] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Invalid status value",
				Code:    http.StatusBadRequest,
			})
			return
		}
		task.Status = models.TaskStatus(req.Status)
	}
	if req.Priority != "" {
		// Валидация приоритета
		validPriority := map[string]bool{
			string(models.PriorityLow):    true,
			string(models.PriorityMedium): true,
			string(models.PriorityHigh):   true,
			string(models.PriorityUrgent): true,
		}
		if !validPriority[req.Priority] {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Error:   "Invalid priority value",
				Code:    http.StatusBadRequest,
			})
			return
		}
		task.Priority = models.TaskPriority(req.Priority)
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	result = h.DB.Save(&task)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to update task",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    task,
	})
}

func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	role, _ := c.Get("role")

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	taskID := c.Param("id")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid task ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid request data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Валидация статуса
	validStatus := map[string]bool{
		string(models.StatusPending):    true,
		string(models.StatusInProgress): true,
		string(models.StatusCompleted):  true,
		string(models.StatusCancelled):  true,
	}
	if !validStatus[req.Status] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid status value",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var result *gorm.DB
	if role == "admin" {
		result = h.DB.Model(&models.Task{}).
			Where("id = ?", taskUUID).
			Update("status", req.Status)
	} else {
		result = h.DB.Model(&models.Task{}).
			Where("id = ? AND created_by = ?", taskUUID, userUUID).
			Update("status", req.Status)
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to update task status",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "Task not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Получаем обновленную задачу
	var task models.Task
	h.DB.Where("id = ?", taskUUID).First(&task)

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    task,
	})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	role, _ := c.Get("role")

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	taskID := c.Param("id")
	taskUUID, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Error:   "Invalid task ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var result *gorm.DB
	if role == "admin" {
		result = h.DB.Where("id = ?", taskUUID).Delete(&models.Task{})
	} else {
		result = h.DB.Where("id = ? AND created_by = ?", taskUUID, userUUID).Delete(&models.Task{})
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error:   "Failed to delete task",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Error:   "Task not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    gin.H{"message": "Task deleted successfully"},
	})
}
