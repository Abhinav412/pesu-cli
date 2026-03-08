package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID       uuid.UUID `gorm:"type:uuid"`
	AssignmentID uuid.UUID `gorm:"type:uuid"`
	Status       string    `gorm:"type:varchar(20);not null;default:'queued';check:status IN ('queued', 'processing', 'passed', 'failed', 'error')"`
	Score        int       `gorm:"default:0"`
	Feedback     string    `gorm:"type:text"`
	SubmittedAt  time.Time `gorm:"autoCreateTime"`
	ProcessedAt  time.Time
}
