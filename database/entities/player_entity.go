package entities

import (
	"github.com/google/uuid"
)

// BuffOrDebuffAlias представляет тип дебаффа
type BuffOrDebuffAlias string

const (
	AliasPoisonReason           BuffOrDebuffAlias = "poison"
	AliasMagicReason            BuffOrDebuffAlias = "magic"
	DebuffAliasMagicOffReason   BuffOrDebuffAlias = "magic_off"
	DebuffAliasSkipReason       BuffOrDebuffAlias = "skip"
	DebuffAliasLuckyStoneReason BuffOrDebuffAlias = "lucky_stone"
)

type Weapons struct {
	Name       string
	Damage     uint
	MinCubeHit uint
	Item       string
	Count      uint `json:"Count,omitempty"`
}

type Meds struct {
	Name  string
	Item  string
	Count int
}

type Debuff struct {
	Alias       BuffOrDebuffAlias
	Health      *uint              `json:"health,omitempty"`
	MinCubeHit  *uint              `json:"min_cube_hit,omitempty"`
	BattleStart *BattleStartOption `json:"battle_start,omitempty"`
	Duration    *uint              `json:"duration,omitempty"`
}

type Buff struct {
	Alias       BuffOrDebuffAlias
	Health      *uint              `json:"health,omitempty"`
	MinCubeHit  *uint              `json:"min_cube_hit,omitempty"`
	BattleStart *BattleStartOption `json:"battle_start,omitempty"`
	Duration    *uint              `json:"duration,omitempty"`
}

type Bag struct {
	Name        string
	Description string
}

type Player struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	SectionID       uuid.UUID `gorm:"type:uuid;not null"`
	Section         Section   `gorm:"foreignKey:SectionID"`
	PlayerSection   []PlayerSection
	Health          uint          `gorm:"type:uint;default=0" json:"health"`
	HealthMax       uint          `gorm:"type:uint;default=0" json:"health_max"`
	ReturnToSection uint          `gorm:"type:uint;default=false" json:"return_to_section"`
	Weapons         []Weapons     `gorm:"type:json;nullable;serializer:json" json:"weapons"`
	Meds            Meds          `gorm:"type:json;nullable;serializer:json" json:"meds"`
	Bag             []Bag         `gorm:"type:json;nullable;serializer:json" json:"bag"`
	Debuff          []Debuff      `gorm:"type:json;nullable;serializer:json" json:"debuff"`
	Buff            []Buff        `gorm:"type:json;nullable;serializer:json" json:"buff"`
	Gold            uint          `gorm:"type:uint;default:0" json:"gold"`
	Bonus           []PlayerBonus `gorm:"type:json;nullable;serializer:json" json:"bonus"`
}
