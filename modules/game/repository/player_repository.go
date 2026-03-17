package repository

import (
	"context"
	"errors"
	"gamebook-backend/database/entities"
	playerDto "gamebook-backend/modules/game/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{
		db: db,
	}
}

func (r *playerRepository) GetByUserId(
	ctx context.Context,
	tx *gorm.DB,
	userId string,
) (entities.Player, error) {
	if tx == nil {
		tx = r.db
	}

	var player entities.Player
	result := tx.WithContext(ctx).
		Preload("Section").
		Preload("Section.Transitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("transitions.text_order ASC")
		}).
		Preload("Section.SectionEnemies").
		Preload("PlayerSection").
		Preload("PlayerSection.Section").
		Where("user_id = ?", userId).
		Take(&player)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entities.Player{}, playerDto.ErrPlayerNotFound
	} else if result.Error != nil {
		return entities.Player{}, result.Error
	}

	return player, nil
}

func (r *playerRepository) GetByPlayerId(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) (entities.Player, error) {
	if tx == nil {
		tx = r.db
	}

	var player entities.Player
	result := tx.WithContext(ctx).
		Preload("Section").
		Preload("Section.Transitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("transitions.text_order ASC")
		}).
		Preload("Section.SectionEnemies").
		Preload("PlayerSection").
		Preload("PlayerSection.Section").
		Where("id = ?", playerID).
		Take(&player)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entities.Player{}, playerDto.ErrPlayerNotFound
	} else if result.Error != nil {
		return entities.Player{}, result.Error
	}

	return player, nil
}

func (r *playerRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	userId uuid.UUID,
	section entities.Section,
) (entities.Player, error) {
	player := entities.Player{
		ID:        uuid.New(),
		UserID:    userId,
		SectionID: section.ID,
		Health:    0,
		HealthMax: 0,
		Meds:      entities.Meds{},
		Weapons:   make([]entities.Weapons, 0),
		Bag:       make([]entities.Bag, 0),
		Bonus:     make([]entities.PlayerBonus, 0),
		Debuff:    make([]entities.Debuff, 0),
		Buff:      make([]entities.Buff, 0),
		Gold:      0,
	}

	db := r.db
	if tx != nil {
		db = tx
	}

	if err := db.WithContext(ctx).Create(&player).Error; err != nil {
		return entities.Player{}, err
	}

	return player, nil
}

func (r *playerRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	player entities.Player,
) (*entities.Player, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	if err := db.WithContext(ctx).Save(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}
