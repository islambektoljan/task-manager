package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title       string    `gorm:"not null"`
	Description string
	Status      string `gorm:"default:'pending'"`
	Priority    string `gorm:"default:'medium'"`
	DueDate     time.Time
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
