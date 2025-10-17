package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserProfile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	FirstName string
	LastName  string
	Bio       string
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (up *UserProfile) BeforeCreate(tx *gorm.DB) error {
	up.ID = uuid.New()
	return nil
}
