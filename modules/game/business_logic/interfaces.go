package business_logic

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
)

type BattleLogic interface {
	ProcessBattle(ctx context.Context, player entities.Player, req dto.ActionRequest) (dto.ActionResponse, error)
}

type DiceLogic interface {
	RollDice(ctx context.Context, player entities.Player) (dto.RollTheDiceDto, error)
}

type BonusLogic interface {
	ProcessBonus(ctx context.Context, player entities.Player, bonusAlias string) (dto.ActionResponse, error)
}

type BribeLogic interface {
	ProcessBribe(ctx context.Context, player entities.Player) (dto.ActionResponse, error)
	GetBribeTransition() dto.TransitionDTO
}

type SleepLogic interface {
	ProcessSleep(ctx context.Context, player entities.Player, dice1, dice2 uint) (dto.SleepDTO, error)
}

type SectionLogic interface {
	CheckTransition(ctx context.Context, transition entities.Transition, dices *entities.Dice, player entities.Player) (bool, error)
	UpdatePlayerSection(ctx context.Context, player entities.Player, transition entities.Transition) error
}
