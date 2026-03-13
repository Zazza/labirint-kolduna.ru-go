package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/battle"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"

	"gorm.io/gorm"
)

type sleep9 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep9(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep9{
		db:     db,
		player: player,
	}
}

func (s *sleep9) Execute(
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
			"<p>Соловей-разбойник убил тебя</p>",
		)
	} else {
		exit = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Ты победил</p>",
		)
	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   death,
		NextTry: false,
	}, nil
}
