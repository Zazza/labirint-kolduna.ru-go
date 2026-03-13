package repository

import (
	"context"
	"gamebook-backend/database/entities"
	enemyDTO "gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sectionRepository struct {
	db *gorm.DB
}

func NewSectionRepository(db *gorm.DB) SectionRepository {
	return &sectionRepository{
		db: db,
	}
}

func (r *sectionRepository) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (entities.Section, error) {
	if tx == nil {
		tx = r.db
	}

	var section entities.Section
	if err := tx.WithContext(ctx).Where("id = ?", id).Take(&section).Error; err != nil {
		return entities.Section{}, err
	}

	return section, nil
}

func (r *sectionRepository) GetBySectionNumber(
	ctx context.Context,
	tx *gorm.DB,
	sectionNumber uint,
) (entities.Section, error) {
	if tx == nil {
		tx = r.db
	}

	var section entities.Section
	if err := tx.WithContext(ctx).Where("number = ?", sectionNumber).Take(&section).Error; err != nil {
		return entities.Section{}, err
	}

	return section, nil
}

func (r *sectionRepository) GetListBySectionNumbers(ctx context.Context, tx *gorm.DB, sections []uint) ([]entities.Section, error) {
	if tx == nil {
		tx = r.db
	}

	var result []entities.Section
	if err := tx.WithContext(ctx).Order("number ASC").Where("number IN ?", sections).Find(&result).Error; err != nil {
		return []entities.Section{}, err
	}

	return result, nil
}

func (r *sectionRepository) IsDicesRequired(section entities.Section) bool {
	if len(section.SectionEnemies) > 0 {
		if section.BattleStart != nil && *section.BattleStart == enemyDTO.BattleStartDices {
			return true
		}
	}

	if section.Dices {
		return true
	}

	return false
}

func (r *sectionRepository) GetAllWithTransitions(ctx context.Context, tx *gorm.DB, sectionNumbers []uint) ([]entities.Section, error) {
	if tx == nil {
		tx = r.db
	}

	var result []entities.Section
	if err := tx.WithContext(ctx).
		Preload("Transitions").
		Preload("Transitions.TargetSection").
		Order("number ASC").
		Where("number IN ?", sectionNumbers).
		Find(&result).Error; err != nil {
		return []entities.Section{}, err
	}

	return result, nil
}
