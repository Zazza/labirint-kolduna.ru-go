package repository

import (
	"context"
	"gamebook-backend/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type battleRepository struct {
	db *gorm.DB
}

func NewBattleRepository(db *gorm.DB) BattleRepository {
	return &battleRepository{
		db: db,
	}
}

func (r *battleRepository) GetByPlayerIdAndSectionNumber(
	tx *gorm.DB,
	playerId uuid.UUID,
	section uint,
) ([]entities.Battle, error) {
	if tx == nil {
		tx = r.db
	}

	var battleLog []entities.Battle
	err := tx.
		Order("created_at ASC").
		Where("player_id = ? AND section = ?", playerId, section).
		Find(&battleLog).
		Error
	if err != nil {
		return nil, err
	}

	return battleLog, nil
}

func (r *battleRepository) FindLastByPlayerIdAndSectionNumber(
	tx *gorm.DB,
	playerID uuid.UUID,
	section uint,
) (*entities.Battle, error) {
	if tx == nil {
		tx = r.db
	}

	var record entities.Battle
	result := tx.
		Limit(1).
		Order("created_at DESC").
		Where("player_id = ? AND section = ?", playerID, section).
		Find(&record)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &record, nil
}

func (r *battleRepository) AddRecord(
	tx *gorm.DB,
	battle entities.Battle,
) (entities.Battle, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Create(&battle).Error; err != nil {
		return entities.Battle{}, err
	}

	return battle, nil
}

func (r *battleRepository) RemoveSleepyByPlayerIDAndSectionNumber(
	ctx context.Context,
	tx *gorm.DB,
	playerId uuid.UUID,
	sectionNumber uint,
) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(
		&entities.Battle{},
		"player_id = ? AND section = ?",
		playerId,
		sectionNumber,
	).Error; err != nil {
		return err
	}

	return nil
}
