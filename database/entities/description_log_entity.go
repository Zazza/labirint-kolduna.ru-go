package entities

import (
	"time"

	"github.com/google/uuid"
)

type DescriptionLog struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PlayerSectionID uuid.UUID     `gorm:"type:uuid;not null"`
	PlayerSection   PlayerSection `gorm:"foreignKey:PlayerSectionID"`
	Description     string        `gorm:"type:text;not null" json:"description"`
	CreatedAt       time.Time     `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
