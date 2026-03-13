package repository

import (
	"context"
	"gamebook-backend/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type bonusRepository struct {
	db *gorm.DB
}

func NewBonusRepository(db *gorm.DB) BonusRepository {
	return &bonusRepository{
		db: db,
	}
}

func (r *bonusRepository) GetByPlayerID(ctx context.Context, playerID uuid.UUID) ([]entities.PlayerBonus, error) {
	var player entities.Player
	result := r.db.WithContext(ctx).Where("id = ?", playerID).Take(&player)

	if result.Error != nil {
		return nil, result.Error
	}

	return player.Bonus, nil
}

func (r *bonusRepository) Create(ctx context.Context, db *gorm.DB, bonus *entities.PlayerBonus) error {
	if db == nil {
		db = r.db
	}

	if err := db.WithContext(ctx).Create(bonus).Error; err != nil {
		return err
	}

	return nil
}

func (r *bonusRepository) Update(ctx context.Context, db *gorm.DB, bonus *entities.PlayerBonus) error {
	if db == nil {
		db = r.db
	}

	if err := db.WithContext(ctx).Save(bonus).Error; err != nil {
		return err
	}

	return nil
}

func (r *bonusRepository) Delete(ctx context.Context, db *gorm.DB, playerID uuid.UUID, alias string) error {
	if db == nil {
		db = r.db
	}

	var player entities.Player
	if err := db.WithContext(ctx).Where("id = ?", playerID).Take(&player).Error; err != nil {
		return err
	}

	var updatedBonuses []entities.PlayerBonus
	for _, b := range player.Bonus {
		if b.Alias == nil || *b.Alias != alias {
			updatedBonuses = append(updatedBonuses, b)
		}
	}

	if err := db.WithContext(ctx).Model(&player).Update("bonus", updatedBonuses).Error; err != nil {
		return err
	}

	return nil
}
