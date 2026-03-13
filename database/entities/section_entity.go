package entities

import (
	"github.com/google/uuid"
)

type SectionType string

type ChoiceItems struct {
	Name        string
	Description string
}

type Choice struct {
	Items         []ChoiceItems
	MaxSelections *uint
}

type BribeSuccess struct {
	//SectionNumber *uint `json:"section_number,omitempty"`
}

type BribeFail struct {
	//SectionNumber *uint `json:"section_number,omitempty"`
}

type Bribe struct {
	Text        *string       `json:"text,omitempty"`
	Amount      *uint         `json:"amount,omitempty"`
	AmountDice  *bool         `json:"amount_dice,omitempty"`
	AmountDices *bool         `json:"amount_dices,omitempty"`
	MinDiceHit  *string       `json:"min_dice_hit,omitempty"`
	MinDicesHit *string       `json:"min_dices_hit,omitempty"`
	Success     *BribeSuccess `json:"success,omitempty"`
	Fail        *BribeFail    `json:"fail,omitempty"`
}

type Section struct {
	ID             uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Type           SectionType `gorm:"type:string;not null;default='normal''" json:"type"`
	Number         uint        `gorm:"type:uint;not null;unique;index" json:"number"`
	Text           string      `gorm:"type:text;not null" json:"text"`
	Transitions    []Transition
	Choice         *Choice   `gorm:"type:json;nullable;serializer:json" json:"choice"`
	SectionEnemies []*Enemy  `gorm:"many2many:section_enemies"`
	BattleStart    *string   `gorm:"type:string;nullable" json:"battle_start"`
	BattleSteps    []*string `gorm:"type:json;nullable;serializer:json" json:"battle_steps"`
	Dices          bool      `gorm:"type:bool;not null;default:false" json:"dices"`
	Bribe          *Bribe    `gorm:"type:json;nullable;serializer:json" json:"bribe"`
}
