package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/battle"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"

	"gorm.io/gorm"
)

type sleep3 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep3(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep3{
		db:     db,
		player: player,
	}
}

func (s sleep3) Execute(
	ctx context.Context,
	_ uint,
	_ uint,
) (dto.SleepyKingdomDTO, error) {
	exit := false
	death := false

	sleepyBattle, err := battle.NewCommon(ctx, s.db, &s.player)
	if err != nil {
		return dto.SleepyKingdomDTO{}, err
	}

	if !sleepyBattle.IsWin() {
		death = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Дрёма убил тебя</p>",
		)
	} else {
		exit = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Ты победил и выходишь из Сонного Царства</p>",
		)
	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   death,
		NextTry: false,
	}, nil
}
