package event

import (
	"gamebook-backend/database/entities"

	"github.com/google/uuid"
)

type EnemyUpdateEvent struct {
	PlayerID    uuid.UUID
	SectionID   uuid.UUID
	EnemyID     uuid.UUID
	Health      *uint
	Debuff      *[]entities.Debuff
	Buff        *[]entities.Buff
	Weapons     *[]entities.EnemyWeapon
	Description string
}

func (e EnemyUpdateEvent) GetName() string {
	return "enemy_update"
}
