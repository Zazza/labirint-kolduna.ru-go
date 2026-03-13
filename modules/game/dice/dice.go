package dice

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"
	"math/rand"

	"gorm.io/gorm"
)

type RollTheDices interface {
	RollTheDice(ctx context.Context, player entities.Player) (*uint, error)
	RollTheDices(ctx context.Context, player entities.Player) (*uint, *uint, error)
	StoreDices(
		ctx context.Context,
		player entities.Player,
		diceFirst, diceSecond uint,
		reason dto.ReasonType,
	) error
}

type rollTheDices struct {
	db             *gorm.DB
	diceRepository repository.DiceRepository
}

func NewRollTheDices(
	db *gorm.DB,
	player *entities.Player,
) RollTheDices {
	if player != nil {
		return NewPlayerRollTheDices(db, player)
	}

	return NewEntityRollTheDices(db)
}

func NewRollTheDicesWithRepository(
	diceRepository repository.DiceRepository,
) RollTheDices {
	return &rollTheDices{
		diceRepository: diceRepository,
	}
}

func (d *rollTheDices) RollTheDice(ctx context.Context, player entities.Player) (*uint, error) {
	diceFirst := uint(rand.Intn(5) + 1)
	return &diceFirst, nil
}

func (d *rollTheDices) RollTheDices(ctx context.Context, player entities.Player) (*uint, *uint, error) {
	diceFirst := uint(rand.Intn(5) + 1)
	diceSecond := uint(rand.Intn(5) + 1)
	return &diceFirst, &diceSecond, nil
}

func (d *rollTheDices) StoreDices(
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

func RandomInt(max int) int {
	return rand.Intn(max)
}
