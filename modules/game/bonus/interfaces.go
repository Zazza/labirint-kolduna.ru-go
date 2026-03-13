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
