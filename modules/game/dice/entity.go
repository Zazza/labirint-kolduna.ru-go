package dice

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"
	"math/rand"

	"gorm.io/gorm"
)

type EntityRollTheDices interface {
	RollTheDices
}

type entityRollTheDices struct {
	db             *gorm.DB
	diceRepository repository.DiceRepository
}

func NewEntityRollTheDices(
	db *gorm.DB,
) RollTheDices {
	diceRepository := repository.NewDiceRepository(db)
	return &entityRollTheDices{
		db:             db,
		diceRepository: diceRepository,
	}
}

func NewEntityRollTheDicesWithRepository(
	diceRepository repository.DiceRepository,
) RollTheDices {
	return &entityRollTheDices{
		diceRepository: diceRepository,
	}
}

func (d *entityRollTheDices) RollTheDice(ctx context.Context, _ entities.Player) (*uint, error) {
	diceFirst := uint(rand.Intn(5) + 1)

	return &diceFirst, nil
}

func (d *entityRollTheDices) RollTheDices(ctx context.Context, _ entities.Player) (*uint, *uint, error) {
	diceFirst := uint(rand.Intn(5) + 1)
	diceSecond := uint(rand.Intn(5) + 1)

	return &diceFirst, &diceSecond, nil
}

func (d *entityRollTheDices) StoreDices(
	ctx context.Context,
	player entities.Player,
	diceFirst, diceSecond uint,
	reason dto.ReasonType,
) error {
	_, err := d.diceRepository.Create(ctx, d.db, player.ID, diceFirst, diceSecond, reason)
	if err != nil {
		return err
	}

	return nil
}
