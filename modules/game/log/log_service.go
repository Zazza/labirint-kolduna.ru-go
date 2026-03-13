package log

import "github.com/google/uuid"

type PlayerLogService interface {
	LogDiceRoll(playerID uuid.UUID, dice1, dice2 uint, reason string)
	LogBattleHit(playerID uuid.UUID, attacker string, defender string, damage uint, weapon string, buffs, debuffs []string)
	LogBribe(playerID uuid.UUID, success bool, amount uint, target string)
	LogBonusUsed(playerID uuid.UUID, bonusName string, alias string, option *string)
	LogTransition(playerID uuid.UUID, fromSection uint, toSection uint, conditions string)
	GetLogs(playerID uuid.UUID) ([]PlayerLogDTO, error)
}

type PlayerLogDTO struct {
	ActionType  string                 `json:"action_type"`
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details"`
	CreatedAt   string                 `json:"created_at"`
}
