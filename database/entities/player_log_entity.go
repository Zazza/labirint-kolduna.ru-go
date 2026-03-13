package entities

import (
	"time"

	"github.com/google/uuid"
)

type PlayerLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PlayerID    uuid.UUID `gorm:"type:uuid;not null;index"`
	ActionType  string    `gorm:"type:varchar(50);not null;index"`
	Description string    `gorm:"type:text;not null"`
	Details     string    `gorm:"type:jsonb;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
