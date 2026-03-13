package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type BattleStartOption string

const (
	BattleStartOptionEnemy  BattleStartOption = "enemy"
	BattleStartOptionPlayer BattleStartOption = "player"
	BattleStartOptionDices  BattleStartOption = "dices"
)

type PlayerChangeWeapon struct {
	Item   *string `json:"Item,omitempty"`
	Change *string `json:"Change,omitempty"`
}

type PlayerChangeMeds struct {
	Item   *string `json:"Item,omitempty"`
	Change *string `json:"Change,omitempty"`
}

type PlayerBonusOption struct {
	Alias *string `json:"alias,omitempty"`
	Name  *string `json:"name,omitempty"`
}

type PlayerBonus struct {
	Alias  *string              `json:"alias,omitempty"`
	Name   *string              `json:"name,omitempty"`
	Option *[]PlayerBonusOption `json:"option,omitempty"`
}

type PlayerChange struct {
	Health          *string               `json:"health,omitempty"`
	Gold            *string               `json:"gold,omitempty"`
	Weapons         *[]PlayerChangeWeapon `json:"weapons,omitempty"`
	Meds            *[]PlayerChangeMeds   `json:"meds,omitempty"`
	Bag             *[]Bag                `json:"bag,omitempty"`
	Bonus           *[]PlayerBonus        `json:"bonus,omitempty"`
	ReturnToSection *uint                 `json:"return_to_section,omitempty"`
}

type PlayerDebuff struct {
	Alias        *BuffOrDebuffAlias `json:"alias,omitempty"`
	Health       *uint              `json:"health,omitempty"`
	HealthString *string            `json:"health_string,omitempty"`
	MinCubeHit   *uint              `json:"min_cube_hit,omitempty"`
	BattleStart  *BattleStartOption `json:"battle_start,omitempty"`
	Duration     *uint              `json:"duration,omitempty"`
}

type PlayerBuff struct {
	Alias        *BuffOrDebuffAlias `json:"alias,omitempty"`
	Health       *uint              `json:"health,omitempty"`
	HealthString *string            `json:"health_string,omitempty"`
	MinCubeHit   *uint              `json:"min_cube_hit,omitempty"`
	BattleStart  *BattleStartOption `json:"battle_start,omitempty"`
	Duration     *uint              `json:"duration,omitempty"`
}

type Transition struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TextOrder       uint            `gorm:"not null;default=0" json:"text_order"`
	SectionID       uuid.UUID       `gorm:"type:uuid;not null"`
	Section         Section         `gorm:"foreignKey:SectionID"`
	TargetSectionID uuid.UUID       `gorm:"type:uuid;not null"`
	TargetSection   Section         `gorm:"foreignKey:TargetSectionID"`
	AvailableOnce   bool            `gorm:"type:bool;not null;default:false" json:"available_once"`
	Text            string          `gorm:"type:text;not null" json:"text"`
	IsBattleWin     *bool           `gorm:"type:bool;nullable" json:"is_battle_win"`
	BribeResult     *bool           `gorm:"type:bool;nullable" json:"bribe_result"`
	PlayerInput     bool            `gorm:"type:bool;default=false" json:"player_input"`
	Dice            *[]string       `gorm:"type:json;nullable;serializer:json" json:"dice"`
	Dices           *[]string       `gorm:"type:json;nullable;serializer:json" json:"dices"`
	Condition       *string         `gorm:"type:text;nullable" json:"condition"`
	PlayerChange    *PlayerChange   `gorm:"column:player_change;type:jsonb;nullable;serializer:json" json:"-"`
	PlayerDebuff    []*PlayerDebuff `gorm:"type:json;nullable;serializer:json" json:"player_debuff"`
	PlayerBuff      []*PlayerBuff   `gorm:"type:json;nullable;serializer:json" json:"player_buff"`
}

// Scan - для чтения из базы данных
func (pc *PlayerChange) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, pc)
}

// Value - для записи в базу данных
func (pc *PlayerChange) Value() (driver.Value, error) {
	if pc == nil {
		return nil, nil
	}
	if pc.Health == nil && pc.Gold == nil && pc.Weapons == nil && pc.Meds == nil && pc.Bag == nil && pc.Bonus == nil && pc.ReturnToSection == nil {
		return nil, nil
	}
	return json.Marshal(pc)
}
