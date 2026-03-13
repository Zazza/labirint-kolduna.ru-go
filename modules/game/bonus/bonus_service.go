package bonus

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/log"
	"gamebook-backend/modules/game/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BonusService interface {
	GetBonus(ctx context.Context, db *gorm.DB, player entities.Player, bonusAlias string) (dto.Bonus, error)
	UseBonus(ctx context.Context, db *gorm.DB, player entities.Player, bonusAlias string) (dto.ActionResponse, error)
	GetAvailableBonuses(ctx context.Context, db *gorm.DB, playerID uuid.UUID) ([]dto.PlayerInfoBonus, error)
}

type bonusService struct {
	bonusRepository repository.BonusRepository
	factory         BonusFactory
	logService      log.PlayerLogService
}

func NewBonusService(
	bonusRepo repository.BonusRepository,
	factory BonusFactory,
) BonusService {
	return &bonusService{
		bonusRepository: bonusRepo,
		factory:         factory,
	}
}

func NewBonusServiceWithLogging(
	bonusRepo repository.BonusRepository,
	factory BonusFactory,
	logService log.PlayerLogService,
) BonusService {
	return &bonusService{
		bonusRepository: bonusRepo,
		factory:         factory,
		logService:      logService,
	}
}

func (s *bonusService) GetBonus(
	ctx context.Context,
	db *gorm.DB,
	player entities.Player,
	bonusAlias string,
) (dto.Bonus, error) {
	return s.factory.GetBonus(db, player, bonusAlias)
}

func (s *bonusService) UseBonus(
	ctx context.Context,
	db *gorm.DB,
	player entities.Player,
	bonusAlias string,
) (dto.ActionResponse, error) {
	bonus, err := s.GetBonus(ctx, db, player, bonusAlias)
	if err != nil {
		return dto.ActionResponse{
			Result: dto.ResultFalse,
			Error:  err.Error(),
		}, err
	}

	req := dto.BonusRequest{
		Bonus: bonusAlias,
	}

	err = bonus.Execute(ctx, req)
	if err != nil {
		return dto.ActionResponse{
			Result: dto.ResultFalse,
			Error:  err.Error(),
		}, err
	}

	if s.logService != nil {
		var option *string
		if req.Option != "" {
			option = &req.Option
		}
		s.logService.LogBonusUsed(player.ID, bonusAlias, bonusAlias, option)
	}

	return dto.ActionResponse{
		Result: dto.ResultTrue,
	}, nil
}

func (s *bonusService) GetAvailableBonuses(
	ctx context.Context,
	db *gorm.DB,
	playerID uuid.UUID,
) ([]dto.PlayerInfoBonus, error) {
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
