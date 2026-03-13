package entities

import (
	"github.com/google/uuid"
)

type EnemyWeapon struct {
	Name       string
	Damage     uint
	MinCubeHit uint
	Item       string
	Count      *int `json:"count,omitempty"`
}

type PlayerSectionEnemy struct {
	ID        uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	PlayerID  uuid.UUID     `gorm:"type:uuid;not null"`
	Player    Player        `gorm:"foreignKey:PlayerID"`
	SectionID uuid.UUID     `gorm:"type:uuid;not null"`
	Section   Section       `gorm:"foreignKey:SectionID"`
	EnemyID   uuid.UUID     `gorm:"type:uuid;not null"`
	Enemy     Enemy         `gorm:"foreignKey:EnemyID"`
	Health    uint          `gorm:"type:uint;not null;" json:"health"`
	Debuff    []Debuff      `gorm:"type:json;nullable;serializer:json" json:"debuff"`
	Buff      []Buff        `gorm:"type:json;nullable;serializer:json" json:"buff"`
	Weapons   []EnemyWeapon `gorm:"type:json;nullable;serializer:json" json:"weapons"`

	Timestamp
}
