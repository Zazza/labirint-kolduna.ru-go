package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"

	"gorm.io/gorm"
)

type sleep4 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep4(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep4{
		db:     db,
		player: player,
	}
}

func (s sleep4) Execute(
	ctx context.Context,
	dice1 uint,
	_ uint,
) (dto.SleepyKingdomDTO, error) {
	exit := false
	death := false

	if dice1 < 5 {
		death = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>1-4 очка - ты съеден</p>",
		)
	}

	if 5 <= dice1 {
		exit = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Если выбросишь 5 или 6 очков, он пройдет мимо</p>",
		)
	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   death,
		NextTry: false,
	}, nil
}
