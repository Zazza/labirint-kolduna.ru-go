package dice

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bonus/helpers"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/repository"
	"math/rand"

	"gorm.io/gorm"
)

type PlayerRollTheDices interface {
	RollTheDices
}

type playerRollTheDices struct {
	db             *gorm.DB
	diceRepository repository.DiceRepository
	player         entities.Player
}

func NewPlayerRollTheDices(
	db *gorm.DB,
	player *entities.Player,
) PlayerRollTheDices {
	diceRepository := repository.NewDiceRepository(db)

	return &playerRollTheDices{
		db:             db,
		diceRepository: diceRepository,
		player:         *player,
	}
}

func (d *playerRollTheDices) RollTheDice(ctx context.Context, player entities.Player) (*uint, error) {
	spellIncrement := 0
	if helpers.HasBuff(player.Buff, entities.DebuffAliasLuckyStoneReason) {
		spellIncrement = helpers.LuckyStoneIncrement
	}

	diceFirst := uint(rand.Intn(5)+1) + uint(spellIncrement)

	return &diceFirst, nil
}

func (d *playerRollTheDices) RollTheDices(ctx context.Context, player entities.Player) (*uint, *uint, error) {
	spellIncrementFirst := 0
	spellIncrementSecond := 0
	if helpers.HasBuff(player.Buff, entities.DebuffAliasLuckyStoneReason) {
		spellIncrementFirst = rand.Intn(1) + 1
		spellIncrementSecond = helpers.LuckyStoneIncrement - spellIncrementFirst
	}

	diceFirst := uint(rand.Intn(5)+1) + uint(spellIncrementFirst)
	diceSecond := uint(rand.Intn(5)+1) + uint(spellIncrementSecond)

	return &diceFirst, &diceSecond, nil
}

func (d *playerRollTheDices) StoreDices(
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
