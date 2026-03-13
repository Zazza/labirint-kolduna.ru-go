package event

import (
	"gamebook-backend/database/entities"

	"github.com/google/uuid"
)

type PlayerUpdateEvent struct {
	PlayerID        uuid.UUID
	SectionID       *uuid.UUID
	ReturnToSection *uuid.UUID
	Health          *uint
	HealthMax       *uint
	Weapons         *[]entities.Weapons
	Meds            *entities.Meds
	Bag             *[]entities.Bag
	Debuff          *[]entities.Debuff
	DebuffList      *[]entities.Debuff
	Buff            *[]entities.Buff
	BuffList        *[]entities.Buff
	Bonus           *[]entities.PlayerBonus
	BonusList       *[]entities.PlayerBonus
	Gold            *uint
}

func (e PlayerUpdateEvent) GetName() string {
	return "player_update"
}
