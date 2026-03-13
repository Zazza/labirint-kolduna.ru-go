package repository

import (
	"context"
	"errors"
	"gamebook-backend/database/entities"
	diceDTO "gamebook-backend/modules/game/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type diceRepository struct {
	db *gorm.DB
}

func NewDiceRepository(db *gorm.DB) DiceRepository {
	return &diceRepository{
		db: db,
	}
}

func (r *diceRepository) GetLastByPlayerId(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
	reason diceDTO.ReasonType,
) (entities.Dice, error) {
	if tx == nil {
		tx = r.db
	}

	var dice entities.Dice
	result := tx.WithContext(ctx).
		Order("created_at desc").
		Limit(1).
		Where("player_id = ? AND reason = ?", playerID, reason).
		Take(&dice)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return entities.Dice{}, diceDTO.MessageDicesNotDefined
	} else if result.Error != nil {
		return entities.Dice{}, result.Error
	}

	return dice, nil
}

func (r *diceRepository) FindBattleDicesByPlayerId(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
) diceDTO.BattleDicesDTO {
	if tx == nil {
		tx = r.db
	}

	battleDicesExists := false
	dices, err := r.GetLastByPlayerId(ctx, tx, playerID, diceDTO.ReasonBattle)
	if err != nil && !errors.Is(err, diceDTO.MessageDicesNotDefined) {
		battleDicesExists = false
	} else if errors.Is(err, diceDTO.MessageDicesNotDefined) {
		battleDicesExists = false
		err = nil
	} else {
		battleDicesExists = true
	}

	return diceDTO.BattleDicesDTO{
		Exists: battleDicesExists,
		Dices:  dices,
		Error:  err,
	}
}

func (r *diceRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	playerID uuid.UUID,
	diceFirst uint,
	diceSecond uint,
	reason diceDTO.ReasonType,
) (entities.Dice, error) {
	if tx == nil {
		tx = r.db
	}

	dice := entities.Dice{
		ID:         uuid.New(),
		PlayerID:   playerID,
		DiceFirst:  diceFirst,
		DiceSecond: diceSecond,
		Reason:     string(reason),
	}

	if err := tx.WithContext(ctx).Create(&dice).Error; err != nil {
		return entities.Dice{}, err
	}

	return dice, nil
}

func (r *diceRepository) Remove(ctx context.Context, tx *gorm.DB, playerID uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(
		&entities.Dice{},
		"player_id = ?",
		playerID,
	).Error; err != nil {
		return err
	}

	return nil
}
