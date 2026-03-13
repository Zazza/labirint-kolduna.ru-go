package repository

import (
	"context"
	"gamebook-backend/database/entities"
	enemyDTO "gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transitionRepository struct {
	db *gorm.DB
}

func NewTransitionRepository(db *gorm.DB) TransitionRepository {
	return &transitionRepository{
		db: db,
	}
}

func (r *transitionRepository) GetByTransitionID(
	ctx context.Context,
	tx *gorm.DB,
	transitionId uuid.UUID,
) (entities.Transition, error) {
	if tx == nil {
		tx = r.db
	}

	var transition entities.Transition
	if err := tx.WithContext(ctx).Preload("Section").Preload("TargetSection").Where("id = ?", transitionId).Take(&transition).Error; err != nil {
		return entities.Transition{}, err
	}

	return transition, nil
}

func (r *transitionRepository) IsDicesRequired(transition entities.Transition) bool {
	if len(transition.Section.SectionEnemies) > 0 {
		if transition.Section.BattleStart != nil && *transition.Section.BattleStart == enemyDTO.BattleStartDices {
			return true
		}
	}

	if transition.Section.Dices {
		return true
	}

	return false
}
