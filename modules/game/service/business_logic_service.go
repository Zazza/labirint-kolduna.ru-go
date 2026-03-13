package service

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bonus"
	"gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessLogicService interface {
	ProcessBonus(ctx context.Context, db *gorm.DB, player entities.Player, bonusAlias string) (dto.ActionResponse, error)
	GetAvailableBonuses(ctx context.Context, db *gorm.DB, playerID string) ([]dto.PlayerInfoBonus, error)
}

type businessLogicService struct {
	bonusService bonus.BonusService
}

func NewBusinessLogicService(
	bonusSvc bonus.BonusService,
) BusinessLogicService {
	return &businessLogicService{
		bonusService: bonusSvc,
	}
}

func (s *businessLogicService) ProcessBonus(
	ctx context.Context,
	db *gorm.DB,
	player entities.Player,
	bonusAlias string,
) (dto.ActionResponse, error) {
	return s.bonusService.UseBonus(ctx, db, player, bonusAlias)
}

func (s *businessLogicService) GetAvailableBonuses(
	ctx context.Context,
	db *gorm.DB,
	playerID string,
) ([]dto.PlayerInfoBonus, error) {
	playerUUID, err := parsePlayerID(playerID)
	if err != nil {
		return nil, err
	}

	return s.bonusService.GetAvailableBonuses(ctx, db, playerUUID)
}

func parsePlayerID(playerID string) (uuid.UUID, error) {
	return uuid.Parse(playerID)
}
