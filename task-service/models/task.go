package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string
type TaskPriority string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusCancelled  TaskStatus = "cancelled"

	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

type Task struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primary_key" json:"id"`
	Title       string       `gorm:"not null" json:"title" binding:"required"`
	Description string       `json:"description"`
	Status      TaskStatus   `gorm:"default:'pending'" json:"status"`
	Priority    TaskPriority `gorm:"default:'medium'" json:"priority"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	CreatedBy   uuid.UUID    `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	// Устанавливаем значения по умолчанию если не заданы
	if t.Status == "" {
		t.Status = StatusPending
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}

	return nil
}

func (Task) TableName() string {
	return "task_schema.tasks"
}
