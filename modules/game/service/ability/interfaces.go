package ability

import (
	"context"

	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
)

type MedsAbility interface {
	Validate(ctx context.Context, player entities.Player) error
	Execute(ctx context.Context, player entities.Player) (dto.MedsDTO, error)
	Result() dto.MedsResultDTO
}

type BonusAbility interface {
	Validate(ctx context.Context, player entities.Player) error
	Execute(ctx context.Context, player entities.Player, req dto.BonusRequest) error
	Result() dto.BonusDTO
}

type SleepAbility interface {
	Validate(ctx context.Context, player entities.Player) error
	Execute(ctx context.Context, player entities.Player) (dto.SleepDTO, error)
	Result() dto.SleepResultDTO
}

type BribeAbility interface {
	Validate(ctx context.Context, player entities.Player) error
	Execute(ctx context.Context, player entities.Player, req dto.BribeRequest) error
	Result() dto.BribeDTO
}

type DiceAbility interface {
	Validate(ctx context.Context, player entities.Player) error
	Execute(ctx context.Context, player entities.Player) (dto.DiceDTO, error)
	Result() dto.DiceResultDTO
}
