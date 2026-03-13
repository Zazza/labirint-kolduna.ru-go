package entities

import (
	"github.com/google/uuid"
)

type Dice struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PlayerID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Reason     string    `gorm:"type:varchar(16)" json:"reason"`
	DiceFirst  uint      `gorm:"type:uint;not null;" json:"dice_first"`
	DiceSecond uint      `gorm:"type:uint;not null;" json:"dice_second"`

	Timestamp
}
