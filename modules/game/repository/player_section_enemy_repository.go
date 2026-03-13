package repository

import (
	"context"
	"gamebook-backend/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type playerSectionEnemyRepository struct {
	db *gorm.DB
}

func NewPlayerSectionEnemyRepository(db *gorm.DB) PlayerSectionEnemyRepository {
	return &playerSectionEnemyRepository{
		db: db,
	}
}

func (r *playerSectionEnemyRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
	sectionID uuid.UUID,
	sectionEnemies []*entities.Enemy,
) ([]entities.PlayerSectionEnemy, error) {
	if tx == nil {
		tx = r.db
	}

	var result []entities.PlayerSectionEnemy

	for _, item := range sectionEnemies {
		enemy := entities.PlayerSectionEnemy{
			ID:        uuid.New(),
			SectionID: sectionID,
			PlayerID:  playerID,
			Health:    item.Health,
			EnemyID:   item.ID,
			Debuff:    make([]entities.Debuff, 0),
			Buff:      make([]entities.Buff, 0),
			Weapons:   item.Weapons,
		}

		result = append(result, enemy)

		if err := tx.WithContext(ctx).Create(&enemy).Error; err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (r *playerSectionEnemyRepository) GetEnemiesByPlayerIDAndSectionID(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
	sectionID uuid.UUID,
) (*[]entities.PlayerSectionEnemy, error) {
	if tx == nil {
		tx = r.db
	}

	var enemies []entities.PlayerSectionEnemy
	result := tx.WithContext(ctx).
		Where("player_id = ? AND section_id = ?", playerID, sectionID).
		Find(&enemies)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &enemies, nil
}

func (r *playerSectionEnemyRepository) GetEnemiesByPlayerIDAndAlias(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
	enemyAlias string,
) (*[]entities.PlayerSectionEnemy, error) {
	if tx == nil {
		tx = r.db
	}

	var enemies []entities.PlayerSectionEnemy
	result := tx.WithContext(ctx).
		Table("player_section_enemies AS pse").
		Joins("left join enemies on enemies.id = pse.enemy_id").
		Where("pse.player_id = ? AND enemies.alias = ?", playerID, enemyAlias).
		Preload("Enemy").
		Find(&enemies)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &enemies, nil
}

func (r *playerSectionEnemyRepository) UpdateAll(
	ctx context.Context,
	tx *gorm.DB,
	enemies []entities.PlayerSectionEnemy,
) error {
	if tx == nil {
		tx = r.db
	}

	for _, enemy := range enemies {
		if err := r.Update(ctx, tx, enemy); err != nil {
			return err
		}
	}

	return nil
}

func (r *playerSectionEnemyRepository) Update(
	ctx context.Context,
	tx *gorm.DB,
	enemy entities.PlayerSectionEnemy,
) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where(
		"enemy_id = ? AND player_id = ? AND section_id = ?",
		enemy.EnemyID,
		enemy.PlayerID,
		enemy.SectionID,
	).
		Select("Health", "Debuff", "Buff", "Weapons").
		Updates(entities.PlayerSectionEnemy{
			Health:  enemy.Health,
			Debuff:  enemy.Debuff,
			Buff:    enemy.Buff,
			Weapons: enemy.Weapons,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r *playerSectionEnemyRepository) RemoveSleepyByPlayerIDAndSectionNumber(ctx context.Context, tx *gorm.DB, playerId uuid.UUID, sectionId uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(
		&entities.PlayerSectionEnemy{},
		"player_id = ? AND section_id = ?",
		playerId,
		sectionId,
	).Error; err != nil {
		return err
	}

	return nil
}
