package repository

import (
	"context"
	"gamebook-backend/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type playerLogRepository struct {
	db *gorm.DB
}

func NewPlayerLogRepository(db *gorm.DB) PlayerLogRepository {
	return &playerLogRepository{
		db: db,
	}
}

func (r *playerLogRepository) Create(ctx context.Context, db *gorm.DB, log *entities.PlayerLog) error {
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Create(log).Error
}

func (r *playerLogRepository) GetByPlayerID(ctx context.Context, db *gorm.DB, playerID uuid.UUID) ([]entities.PlayerLog, error) {
	if db == nil {
		db = r.db
	}
	var logs []entities.PlayerLog
	err := db.WithContext(ctx).Where("player_id = ?", playerID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}
