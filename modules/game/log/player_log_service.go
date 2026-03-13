package log

import (
	"gamebook-backend/modules/game/channel"
	"gamebook-backend/modules/game/listener/event"

	"github.com/google/uuid"
)

type playerLogService struct{}

func NewPlayerLogService() PlayerLogService {
	return &playerLogService{}
}

func (s *playerLogService) LogDiceRoll(playerID uuid.UUID, dice1, dice2 uint, reason string) {
	var result uint
	if dice2 > 0 {
		result = dice1 + dice2
	} else {
		result = dice1
	}

	description := "Бросок кубика"
	if reason != "" {
		description = "Бросок кубика: " + reason
	}

	details := map[string]interface{}{
		"reason": reason,
		"dice1":  dice1,
		"dice2":  dice2,
		"result": result,
	}

	channel.ChPlayerLog <- event.PlayerLogEvent{
		PlayerID:    playerID,
		ActionType:  "dice_roll",
		Description: description,
		Details:     details,
	}
}

func (s *playerLogService) LogBattleHit(playerID uuid.UUID, attacker string, defender string, damage uint, weapon string, buffs, debuffs []string) {
	description := "Удар в бою"
	if attacker == "player" {
		description = "Ты нанес удар"
	} else {
		description = "Враг нанес удар"
	}

	details := map[string]interface{}{
		"attacker": attacker,
		"defender": defender,
		"damage":   damage,
		"weapon":   weapon,
	}
	if len(buffs) > 0 {
		details["buffs"] = buffs
	}
	if len(debuffs) > 0 {
		details["debuffs"] = debuffs
	}

	channel.ChPlayerLog <- event.PlayerLogEvent{
		PlayerID:    playerID,
		ActionType:  "battle_hit",
		Description: description,
		Details:     details,
	}
}

func (s *playerLogService) LogBribe(playerID uuid.UUID, success bool, amount uint, target string) {
	description := "Попытка взятки"
	if success {
		description = "Успешная взятка"
	} else {
		description = "Неудачная взятка"
	}

	details := map[string]interface{}{
		"success": success,
		"amount":  amount,
		"target":  target,
	}

	channel.ChPlayerLog <- event.PlayerLogEvent{
		PlayerID:    playerID,
		ActionType:  "bribe",
		Description: description,
		Details:     details,
	}
}

func (s *playerLogService) LogBonusUsed(playerID uuid.UUID, bonusName string, alias string, option *string) {
	description := "Использован бонус: " + bonusName
	if option != nil {
		description += " (" + *option + ")"
	}

	details := map[string]interface{}{
		"bonus_name": bonusName,
		"alias":      alias,
	}
	if option != nil {
		details["option"] = *option
	}

	channel.ChPlayerLog <- event.PlayerLogEvent{
		PlayerID:    playerID,
		ActionType:  "bonus_used",
		Description: description,
		Details:     details,
	}
}

func (s *playerLogService) LogTransition(playerID uuid.UUID, fromSection uint, toSection uint, conditions string) {
	description := "Переход к секции " + string(rune(toSection))

	details := map[string]interface{}{
		"from_section": fromSection,
		"to_section":   toSection,
	}
	if conditions != "" {
		details["conditions"] = conditions
	}

	channel.ChPlayerLog <- event.PlayerLogEvent{
		PlayerID:    playerID,
		ActionType:  "transition",
		Description: description,
		Details:     details,
	}
}

func (s *playerLogService) GetLogs(playerID uuid.UUID) ([]PlayerLogDTO, error) {
	return nil, nil
}
