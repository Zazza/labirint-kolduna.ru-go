package battle

import (
	"gamebook-backend/database/entities"
	battleDTO "gamebook-backend/modules/game/dto"
)

type BattleStartStrategy interface {
	IsMyMove(player *entities.Player, dices *entities.Dice, battleLog *entities.Battle, sectionBattleStart *entities.BattleStartOption) (bool, error)
}

type StrategyFactory interface {
	GetStrategy(battleStartOption *entities.BattleStartOption) BattleStartStrategy
}

type battleStrategyFactory struct {
	strategies map[entities.BattleStartOption]BattleStartStrategy
}

func NewBattleStrategyFactory() StrategyFactory {
	factory := &battleStrategyFactory{
		strategies: make(map[entities.BattleStartOption]BattleStartStrategy),
	}

	factory.registerDefaults()
	return factory
}

func (f *battleStrategyFactory) registerDefaults() {
	enemyStart := entities.BattleStartOptionEnemy
	f.strategies[enemyStart] = &EnemyStartStrategy{}

	playerStart := entities.BattleStartOptionPlayer
	f.strategies[playerStart] = &PlayerStartStrategy{}

	diceStart := entities.BattleStartOptionDices
	f.strategies[diceStart] = &DiceStartStrategy{}
}

func (f *battleStrategyFactory) GetStrategy(battleStartOption *entities.BattleStartOption) BattleStartStrategy {
	if battleStartOption == nil {
		return &DefaultBattleStartStrategy{}
	}

	if strategy, exists := f.strategies[*battleStartOption]; exists {
		return strategy
	}

	return &DefaultBattleStartStrategy{}
}

type EnemyStartStrategy struct{}

func (s *EnemyStartStrategy) IsMyMove(player *entities.Player, dices *entities.Dice, battleLog *entities.Battle, sectionBattleStart *entities.BattleStartOption) (bool, error) {
	return false, nil
}

type PlayerStartStrategy struct{}

func (s *PlayerStartStrategy) IsMyMove(player *entities.Player, dices *entities.Dice, battleLog *entities.Battle, sectionBattleStart *entities.BattleStartOption) (bool, error) {
	return true, nil
}

type DiceStartStrategy struct{}

func (s *DiceStartStrategy) IsMyMove(player *entities.Player, dices *entities.Dice, battleLog *entities.Battle, sectionBattleStart *entities.BattleStartOption) (bool, error) {
	if dices == nil {
		return false, battleDTO.MessageDicesNotDefined
	}
	if dices.DiceFirst > dices.DiceSecond {
		return true, nil
	}
	return false, nil
}

type DefaultBattleStartStrategy struct{}

func (s *DefaultBattleStartStrategy) IsMyMove(player *entities.Player, dices *entities.Dice, battleLog *entities.Battle, sectionBattleStart *entities.BattleStartOption) (bool, error) {
	if sectionBattleStart == nil {
		if *player.Section.BattleStart == battleDTO.BattleStartPlayer ||
			(dices != nil && dices.DiceFirst > dices.DiceSecond) ||
			(*player.Section.BattleStart == battleDTO.BattleStartDice1 && dices.DiceFirst > 1) ||
			(*player.Section.BattleStart == battleDTO.BattleStartDice2 && dices.DiceFirst > 2) ||
			(*player.Section.BattleStart == battleDTO.BattleStartDice3 && dices.DiceFirst > 3) ||
			(*player.Section.BattleStart == battleDTO.BattleStartDice4 && dices.DiceFirst > 4) ||
			(*player.Section.BattleStart == battleDTO.BattleStartDice5 && dices.DiceFirst > 5) {
			return true, nil
		}
	}

	return false, nil
}
