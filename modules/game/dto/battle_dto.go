package dto

import (
	"gamebook-backend/database/entities"
	apperrors "gamebook-backend/pkg/errors"
)

const (
	Hand          = "Руки"
	Sword         = "Меч-кладенец"
	Lightning     = "Молнии"
	BallLightning = "Шаровые молнии"
	StepTypeSpell = "spell"
)

type BattleDto struct {
	BattleLog  *entities.Battle
	Enemy      *EnemyDTO
	Step       uint
	PlayerTurn bool
	Finish     bool
}

type BattleResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type EnemyDTO struct {
	Abstract  *entities.Enemy
	Instance  *entities.PlayerSectionEnemy
	Companion *EnemyDTO
}

type EnemyAbstractDTO struct {
	Name         string
	Health       *uint
	MinCubeHit   *uint
	OnlyDiceHits *uint
	DamageType   *string
}

type EnemyInstanceDTO struct {
	EnemyID uint
	Health  uint
	Debuff  *uint
}

type WeaponDamageDto struct {
	Damage      uint
	Description string
}

type BattleState struct {
	Step       uint
	Enemy      *EnemyDTO
	PlayerTurn bool
}

func NewBattleResponse(message string) *BattleResponse {
	return &BattleResponse{
		Success: true,
		Message: message,
	}
}

func NewBattleErrorResponse(message string, err error) *BattleResponse {
	errorCode := ""
	if err != nil {
		errorCode = apperrors.CodeInternalError
	}

	return &BattleResponse{
		Success: false,
		Message: message + ": " + errorCode,
	}
}

func NewBattleResponseWithData(data interface{}) *BattleResponse {
	return &BattleResponse{
		Success: true,
		Message: "",
		Data:    data,
	}
}
