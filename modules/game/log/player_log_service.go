package log

import (
	"context"
	"fmt"
	"gamebook-backend/modules/game/helper"
	"github.com/google/uuid"
)

type playerLogService struct{}

func NewPlayerLogService() PlayerLogService {
	return &playerLogService{}
}

func (s *playerLogService) LogDiceRoll(playerID uuid.UUID, dice1, dice2 uint, reason string) {
	description := "Бросок кубика"
	if reason != "" {
		description = "Бросок кубика: " + reason
	}

	helper.DescriptionMessageWithContext(context.Background(), playerID, fmt.Sprintf("<p>%s</p>", description))
}

func (s *playerLogService) LogBattleHit(playerID uuid.UUID, attacker string, defender string, damage uint, weapon string, buffs, debuffs []string) {
	description := "Удар в бою"
	if attacker == "player" {
		description = "Ты нанес удар"
	} else {
		description = "Враг нанес удар"
	}

	helper.DescriptionMessageWithContext(context.Background(), playerID, fmt.Sprintf("<p>%s</p>", description))
}

func (s *playerLogService) LogBribe(playerID uuid.UUID, success bool, amount uint, target string) {
	description := "Попытка взятки"
	if success {
		description = "Успешная взятка"
	} else {
		description = "Неудачная взятка"
	}

	helper.DescriptionMessageWithContext(context.Background(), playerID, fmt.Sprintf("<p>%s</p>", description))
}

func (s *playerLogService) LogBonusUsed(playerID uuid.UUID, bonusName string, alias string, option *string) {
	description := "Использован бонус: " + bonusName
	if option != nil {
		description += " (" + *option + ")"
	}

	helper.DescriptionMessageWithContext(context.Background(), playerID, fmt.Sprintf("<p>%s</p>", description))
}

func (s *playerLogService) LogTransition(playerID uuid.UUID, fromSection uint, toSection uint, conditions string) {
	description := "Переход к секции " + string(rune(toSection))

	helper.DescriptionMessageWithContext(context.Background(), playerID, fmt.Sprintf("<p>%s</p>", description))
}

func (s *playerLogService) GetLogs(playerID uuid.UUID) ([]PlayerLogDTO, error) {
	return nil, nil
}
