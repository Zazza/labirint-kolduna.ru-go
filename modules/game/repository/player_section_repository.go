package repository

import (
	"context"
	"gamebook-backend/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type playerSectionRepository struct {
	db *gorm.DB
}

func NewPlayerSectionRepository(db *gorm.DB) PlayerSectionRepository {
	return &playerSectionRepository{
		db: db,
	}
}

func (r *playerSectionRepository) Create(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, sectionID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	playerSection := entities.PlayerSection{
		ID:        uuid.New(),
		SectionID: sectionID,
		PlayerID:  playerID,
	}

	if err := tx.WithContext(ctx).Create(&playerSection).Error; err != nil {
		return err
	}

	return nil
}

func (r *playerSectionRepository) UpdateLastTargetSection(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
	targetSectionID uuid.UUID,
) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).
		Order("created_at ASC").
		Limit(1).
		Where("player_id = ?", playerID).
		Select("TargetSectionID").
		Updates(entities.PlayerSection{
			TargetSectionID: targetSectionID,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r *playerSectionRepository) AddDescriptionLog(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
	description string,
) error {
	if tx == nil {
		tx = r.db
	}

	playerSection, err := r.GetLastPlayerSection(ctx, tx, playerID)
	if err != nil {
		return err
	}

	descriptionLog := entities.DescriptionLog{
		ID:              uuid.New(),
		PlayerSectionID: playerSection.ID,
		Description:     description,
	}

	if err := tx.WithContext(ctx).
		Create(&descriptionLog).Error; err != nil {
		return err
	}

	return nil
}

func (r *playerSectionRepository) GetLastSectionDescriptions(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
) (*[]entities.DescriptionLog, error) {
	if tx == nil {
		tx = r.db
	}

	playerSection, err := r.GetLastPlayerSection(ctx, tx, playerID)
	if err != nil {
		return nil, err
	}

	var result []entities.DescriptionLog

	if err := tx.WithContext(ctx).
		Order("created_at ASC").
		Where("player_section_id = ?", playerSection.ID).
		Find(&result).
		Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *playerSectionRepository) GetLastPlayerSection(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) (*entities.PlayerSection, error) {
	if tx == nil {
		tx = r.db
	}

	var playerSection entities.PlayerSection

	result := tx.
		Order("created_at DESC").
		Limit(1).
		WithContext(ctx).
		Preload("Section").
		Where("player_id = ?", playerID).
		Take(&playerSection)
	if result.Error != nil {
		return nil, result.Error
	}

	return &playerSection, nil
}

func (r *playerSectionRepository) RemoveLastPlayerSection(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Select("PlayerSections, DescriptionLogs").Delete(
		&entities.PlayerSection{},
		"player_id IN (SELECT id FROM player_sections AS ps WHERE ps.player_id = ? ORDER BY ps.created_at DESC LIMIT 1)",
		playerID,
	).Error; err != nil {
		return err
	}

	return nil
}

func (r *playerSectionRepository) GetPreviousSectionIdBySectionId(ctx context.Context, tx *gorm.DB, playerID uuid.UUID, sectionID uuid.UUID) (*entities.PlayerSection, error) {
	if tx == nil {
		tx = r.db
	}

	var deathSection entities.PlayerSection
	var playerSection entities.PlayerSection

	result := tx.
		Order("created_at DESC").
		Limit(1).
		WithContext(ctx).
		Where("player_id = ? AND section_id = ?", playerID, sectionID).
		Take(&deathSection)
	if result.Error != nil {
		return nil, result.Error
	}

	result = tx.
		Order("created_at DESC").
		Limit(1).
		WithContext(ctx).
		Preload("Section").
		Where("player_id = ? AND created_at < ?", deathSection.PlayerID, deathSection.CreatedAt).
		Take(&playerSection)
	if result.Error != nil {
		return nil, result.Error
	}

	return &playerSection, nil
}
