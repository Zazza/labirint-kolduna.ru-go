package bonus

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
)

type Bonus interface {
	Use(ctx context.Context, player entities.Player) (dto.ActionResponse, error)
	GetName() string
	GetAlias() string
}

type BonusProcessor interface {
	GetBonus(ctx context.Context, player entities.Player, bonusAlias string) (dto.PlayerInfoBonus, error)
	ProcessBonus(ctx context.Context, player entities.Player, bonusRequest dto.BonusRequest) error
	ValidateBonusAvailability(ctx context.Context, player entities.Player, bonusAlias string) error
	GetAvailableBonuses(ctx context.Context, player entities.Player) ([]dto.PlayerInfoBonus, error)
}

type BonusProvider interface {
	GetBonusByAlias(ctx context.Context, alias string) (dto.PlayerInfoBonus, error)
	ExecuteBonus(ctx context.Context, player entities.Player, bonusAlias string) error
}

type BonusAvailabilityValidator interface {
	ValidateAvailability(ctx context.Context, player entities.Player, bonusAlias string) error
}
