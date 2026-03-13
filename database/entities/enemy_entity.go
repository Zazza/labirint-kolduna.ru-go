package entities

import (
	"github.com/google/uuid"
)

type MagicHit struct {
	Periodicity *uint   `json:"periodicity,omitempty"`
	DicesValues *string `json:"dices_values,omitempty"`
	InstantKill *bool   `json:"instant_kill,omitempty"`
	Damage      *uint   `json:"damage,omitempty"`
	MinDiceHits *string `json:"min_dice_hits,omitempty"`
}

type Enemy struct {
	ID           uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Alias        string            `gorm:"type:varchar(32);not null;unique" json:"alias"`
	Name         string            `gorm:"type:varchar(100);not null" json:"name"`
	Damage       uint              `gorm:"type:uint;not null;default=0" json:"damage"`
	DamageType   BuffOrDebuffAlias `gorm:"type:varchar(32);not null;default=normal" json:"damage_type"`
	MinDiceHits  uint              `gorm:"type:uint;not null;default=6" json:"min_dice_hits"`
	Health       uint              `gorm:"type:uint;not null;default=0" json:"health"`
	Defence      uint              `gorm:"type:uint;not null;default=0" json:"defence"`
	PlayerArmor  bool              `gorm:"type:bool;not null;default=true" json:"player_armor"`
	OnlyDiceHits *uint             `gorm:"type:uint;nullable" json:"only_dice_hits"`
	MagicHit     *MagicHit         `gorm:"type:json;nullable;serializer:json" json:"magic_hit"`
	Weapons      []EnemyWeapon     `gorm:"type:json;nullable;serializer:json" json:"weapons"`
}
