package entities

import (
	"time"

	"github.com/google/uuid"
)

type PlayerSection struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PlayerID        uuid.UUID        `gorm:"type:uuid;not null"`
	Player          Player           `gorm:"foreignKey:PlayerID"`
	SectionID       uuid.UUID        `gorm:"type:uuid;not null"`
	Section         Section          `gorm:"foreignKey:SectionID"`
	TargetSectionID uuid.UUID        `gorm:"type:uuid;nullable"`
	TargetSection   Section          `gorm:"foreignKey:SectionID"`
	Description     []DescriptionLog `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt       time.Time        `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
