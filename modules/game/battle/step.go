package battle

import (
	"gamebook-backend/database/entities"
	battleDTO "gamebook-backend/modules/game/dto"
	"slices"
)

type Step interface {
	GetCurrentBattleSteps() ([]string, error)
	IsMyMove() (bool, error)
	GetNextStepIndex(lastStepIndex uint) uint
	GetCurrentStepIndex() uint
}

type step struct {
	player        *entities.Player
	lastBattleLog *entities.Battle
	battleLog     *[]entities.Battle
	dices         *entities.Dice

	common Common
}

func NewStep(common Common) Step {
	return &step{
		player:        common.GetPlayer(),
		lastBattleLog: common.GetLastBattleLog(),
		battleLog:     common.GetBattleLog(),
		common:        common,
	}
}

func (s *step) GetCurrentBattleSteps() ([]string, error) {
	if s.player.Section.BattleSteps == nil {
		return nil, battleDTO.ErrStepNotDefined
	}

	livingEnemies := uint(0)
	for _, currentEnemy := range *s.common.GetEnemies() {
		if currentEnemy.Health > 0 {
			livingEnemies++
		}
	}

	var result []string
	for i := len(s.player.Section.BattleSteps) - 1; i >= 0; i-- {
		if *s.player.Section.BattleSteps[i] != AttackingPlayer && livingEnemies > 0 {
			result = append(result, *s.player.Section.BattleSteps[i])
			livingEnemies--
		} else if *s.player.Section.BattleSteps[i] == AttackingPlayer {
			result = append(result, *s.player.Section.BattleSteps[i])
		}
	}

	slices.Reverse(result)

	return result, nil
}

func (s *step) IsMyMove() (bool, error) {
	if s.player.Section.BattleSteps == nil {
		return false, battleDTO.ErrStepNotDefined
	}

	for _, debuff := range s.player.Debuff {
		if debuff.BattleStart == nil {
			continue
		}
		switch *debuff.BattleStart {
		case entities.BattleStartOptionEnemy:
			return false, nil
		case entities.BattleStartOptionPlayer:
			return true, nil
		case entities.BattleStartOptionDices:
			if s.dices == nil {
				return false, battleDTO.MessageDicesNotDefined
			}
			if s.dices.DiceFirst > s.dices.DiceSecond {
				return true, nil
			}

			return false, nil
		}
	}

	if s.battleLog == nil {
		if *s.player.Section.BattleStart == battleDTO.BattleStartPlayer ||
			(s.dices.DiceFirst > s.dices.DiceSecond) ||
			(*s.player.Section.BattleStart == battleDTO.BattleStartDice1 && s.dices.DiceFirst > 1) ||
			(*s.player.Section.BattleStart == battleDTO.BattleStartDice2 && s.dices.DiceFirst > 2) ||
			(*s.player.Section.BattleStart == battleDTO.BattleStartDice3 && s.dices.DiceFirst > 3) ||
			(*s.player.Section.BattleStart == battleDTO.BattleStartDice4 && s.dices.DiceFirst > 4) ||
			(*s.player.Section.BattleStart == battleDTO.BattleStartDice5 && s.dices.DiceFirst > 5) {
			return true, nil
		}

		return false, nil
	}

	nextStepIndex := uint(0)
	if s.lastBattleLog != nil {
		nextStepIndex = s.GetNextStepIndex(s.lastBattleLog.Step)
	}

	steps, err := s.GetCurrentBattleSteps()
	if err != nil {
		return false, err
	}

	// TODO: LLM предлагает if len(s.player.Section.BattleSteps) > int(nextStepIndex) {
	//if len(s.player.Section.BattleSteps) <= int(nextStepIndex+1) {
	if len(steps) > int(nextStepIndex) {
		nextStepCharacter := steps[nextStepIndex]
		if nextStepCharacter == AttackingPlayer {
			return true, nil
		}
	} else {
		nextStepCharacter := steps[0]
		if nextStepCharacter == AttackingPlayer {
			return true, nil
		}
	}

	return false, nil
}

func (s *step) GetNextStepIndex(lastStepIndex uint) uint {
	steps, err := s.GetCurrentBattleSteps()
	if err != nil {
		return 0
	}

	if uint(len(steps)) <= lastStepIndex+1 {
		return 0
	}

	return lastStepIndex + 1
}

func (s *step) GetCurrentStepIndex() uint {
	if s.lastBattleLog == nil {
		return 0
	}
	return s.lastBattleLog.Step
}
