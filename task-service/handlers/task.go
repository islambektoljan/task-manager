package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	// Создание задачи
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	// Получение списка задач
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	// Получение задачи по ID
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	// Обновление задачи
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// Удаление задачи
}
