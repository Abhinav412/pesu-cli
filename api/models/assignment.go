package models

import (
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Language    string    `gorm:"type:varchar(50);not null;check:language IN ('python', 'c')"`
	DueDate     time.Time
	CreatedBy   uuid.UUID `gorm:"type:uuid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
