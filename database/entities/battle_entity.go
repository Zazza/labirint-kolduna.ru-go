package entities

import (
	"github.com/google/uuid"
)

type Battle struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PlayerID    uuid.UUID `gorm:"type:uuid;not null;index" json:"player_id"`
	Section     uint      `gorm:"type:uint;not null;index" json:"section"`
	Type        string    `gorm:"type:varchar(16);not null;default=normal" json:"type"`
	Step        uint      `gorm:"type:uint;not null;index" json:"step"`
	Attacking   string    `gorm:"type:varchar(16);not null" json:"attacking"`
	Dice1       uint      `gorm:"type:int;not null" json:"dice1"`
	Dice2       uint      `gorm:"type:int;not null" json:"dice2"`
	Damage      uint      `gorm:"type:int;not null" json:"damage"`
	Description string    `gorm:"type:varchar(1024);nullable" json:"description"`
	Weapon      string    `gorm:"type:varchar(32);not null" json:"weapon"`

	Timestamp
}
