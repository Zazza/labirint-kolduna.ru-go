package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"

	"gorm.io/gorm"
)

type sleep8 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep8(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep8{
		db:     db,
		player: player,
	}
}

func (s *sleep8) Execute(
	ctx context.Context,
	_ uint,
	_ uint,
) (dto.SleepyKingdomDTO, error) {
	rollTheDices := dice.NewRollTheDices(s.db, &s.player)

	diceYou1, diceYou2, err := rollTheDices.RollTheDices(ctx, s.player)
	if err != nil {
		return dto.SleepyKingdomDTO{}, err
	}

	diceEnemy1, diceEnemy2, err := rollTheDices.RollTheDices(ctx, s.player)
	if err != nil {
		return dto.SleepyKingdomDTO{}, err
	}

	exit := false
	death := false
	if *diceYou1+*diceYou2 >= *diceEnemy1+*diceEnemy2 {
		exit = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Ты убежал</p>",
		)
	} else {
		death = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Лихо Одноглазое поймало тебя!</p>",
		)
	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   death,
		NextTry: false,
	}, nil
}
