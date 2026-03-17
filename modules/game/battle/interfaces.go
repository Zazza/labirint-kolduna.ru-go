package battle

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
)

type BattleManager interface {
	StartBattle(ctx context.Context, player entities.Player, section entities.Section) (*entities.Battle, error)
	GetBattle(ctx context.Context, battleID string) (entities.Battle, error)
	GetCurrentBattle(ctx context.Context, player entities.Player) (*entities.Battle, error)
	UpdateBattle(ctx context.Context, battle entities.Battle) error
	EndBattle(ctx context.Context, battle entities.Battle) error
	RollBattleDice(ctx context.Context, player entities.Player) (dto.RollTheDiceDto, error)
	ExecuteBattleStep(ctx context.Context, player entities.Player, step string) (dto.ActionResponse, error)
	GetBattleState(ctx context.Context, player entities.Player) (dto.BattleState, error)
}

type BattleStateProvider interface {
	GetBattleState(ctx context.Context, player entities.Player) (dto.BattleState, error)
}

type DiceRoller interface {
	RollDice(ctx context.Context, player entities.Player) (uint, uint, error)
}

type StepExecutor interface {
	ExecuteStep(ctx context.Context, player entities.Player, step string) (dto.ActionResponse, error)
}
