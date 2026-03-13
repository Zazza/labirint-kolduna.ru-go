package dto

import "gamebook-backend/database/entities"

var (
// ErrRefreshTokenNotFound = errors.Action("refresh token not found")
)

const (
	Hand          = "hand"
	Sword         = "sword"
	Lightning     = "lightning"
	BallLightning = "ball lightning"

	StepTypeNormal = "normal"
	StepTypeSpell  = "spell"
)

type (
	EnemyDTO struct {
		Abstract  *entities.Enemy
		Instance  *entities.PlayerSectionEnemy
		Companion *EnemyDTO
	}

	WeaponDamageDto struct {
		Damage      uint
		Description string
	}

	BattleDto struct {
		Finish bool
		Player *entities.Player
	}
)
