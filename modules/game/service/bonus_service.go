package service

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BonusService interface {
	GetPlayerBonuses(ctx context.Context, db *gorm.DB, playerID uuid.UUID) ([]dto.PlayerInfoBonus, error)
	ConvertBonusesToDTO(bonuses []entities.PlayerBonus) []dto.PlayerInfoBonus
}

type bonusService struct {
	bonusRepository repository.BonusRepository
}

func NewBonusService(
	bonusRepo repository.BonusRepository,
) BonusService {
	return &bonusService{
		bonusRepository: bonusRepo,
	}
}

func (s *bonusService) GetPlayerBonuses(ctx context.Context, db *gorm.DB, playerID uuid.UUID) ([]dto.PlayerInfoBonus, error) {
	bonuses, err := s.bonusRepository.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	return s.ConvertBonusesToDTO(bonuses), nil
}

func (s *bonusService) ConvertBonusesToDTO(bonuses []entities.PlayerBonus) []dto.PlayerInfoBonus {
	var result []dto.PlayerInfoBonus
	for _, playerBonus := range bonuses {
		if playerBonus.Option != nil {
			for _, option := range *playerBonus.Option {
				result = append(result, dto.PlayerInfoBonus{
					Alias:  *playerBonus.Alias,
					Name:   *playerBonus.Name,
					Option: &option,
				})
			}
		} else {
			result = append(result, dto.PlayerInfoBonus{
				Alias: *playerBonus.Alias,
				Name:  *playerBonus.Name,
			})
		}
	}

	return result
}
