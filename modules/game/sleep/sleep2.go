package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"

	"gorm.io/gorm"
)

type sleep2 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep2(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep2{
		db:     db,
		player: player,
	}
}

func (s sleep2) Execute(
	ctx context.Context,
	dice1 uint,
	_ uint,
) (dto.SleepyKingdomDTO, error) {
	var exit bool
	var nextTry bool

	if dice1 <= 3 {
		exit = true
		nextTry = false

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Набрал 1-3 очка - успел обезвредить Дрему</p>",
		)
	}

	if dice1 > 3 {
		exit = false
		nextTry = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>4-6 очков - Дрема усыпил тебя. Это означает, что ты должен входить в cон снова и бросать кубики, чтобы узнать, попадешь ли в Сонное Царство</p>",
		)
	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   false,
		NextTry: nextTry,
	}, nil
}
