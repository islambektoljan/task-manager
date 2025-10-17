package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskSubmission struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	TaskID      uuid.UUID `gorm:"type:uuid;not null"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Solution    string    `gorm:"type:text;not null"`
	Status      string    `gorm:"default:'submitted'"`
	SubmittedAt time.Time
	ReviewedAt  time.Time
	Score       int
	Comments    string
}

func (ts *TaskSubmission) BeforeCreate(tx *gorm.DB) error {
	ts.ID = uuid.New()
	return nil
}
