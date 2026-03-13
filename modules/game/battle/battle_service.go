package battle

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type BattleService interface {
	NewCommon(ctx context.Context, db *gorm.DB, player *entities.Player) (Common, error)
}

type battleService struct {
	battlesRepository      repository.BattleRepository
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository
}

func NewBattleService(
	battlesRepo repository.BattleRepository,
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository,
) BattleService {
	return &battleService{
		battlesRepository:      battlesRepo,
		playerSectionEnemyRepo: playerSectionEnemyRepo,
	}
}

func (s *battleService) NewCommon(
	ctx context.Context,
	db *gorm.DB,
	player *entities.Player,
) (Common, error) {
	return NewCommon(ctx, db, player)
}
