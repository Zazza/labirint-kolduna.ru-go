package sleep

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"

	"gorm.io/gorm"
)

type sleep5 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep5(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep5{
		db:     db,
		player: player,
	}
}

func (s sleep5) Execute(
	ctx context.Context,
	dice1 uint,
	dice2 uint,
) (dto.SleepyKingdomDTO, error) {
	death := false
	exit := false

	if s.player.Health <= dice1+dice2 {
		death = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Пиявки закусали тебя до смерти</p>",
		)
	} else {
		health := s.player.Health - (dice1 + dice2)

		helper.DescriptionMessage(
			s.player.ID,
			fmt.Sprintf("<p>Пиявки покусали тебя и отняли %d HP</p>", health),
		)

		playerUpdateListener, err := listener.HandleEvent(s.db, "player_update")
		if err != nil {
			return dto.SleepyKingdomDTO{}, err
		}
		err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
			PlayerID: s.player.ID,
			Health:   &health,
		})
		if err != nil {
			return dto.SleepyKingdomDTO{}, err
		}

		exit = true
	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   death,
		NextTry: false,
	}, nil
}
