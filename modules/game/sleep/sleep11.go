package sleep

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"

	"gorm.io/gorm"
)

type sleep11 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep11(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep11{
		db:     db,
		player: player,
	}
}

func (s *sleep11) Execute(
	ctx context.Context,
	dice1 uint,
	dice2 uint,
) (dto.SleepyKingdomDTO, error) {
	death := false

	rollTheDices := dice.NewRollTheDices(s.db, &s.player)

	diceDay1, diceDay2, err := rollTheDices.RollTheDices(ctx, s.player)
	if err != nil {
		return dto.SleepyKingdomDTO{}, err
	}

	if s.player.Health <= dice1+dice2 {
		death = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Ты мертв</p>",
		)
	} else {
		health := s.player.Health - (*diceDay1 + *diceDay2)

		helper.DescriptionMessage(
			s.player.ID,
			fmt.Sprintf(
				"Ты провел в темнице %d дней и потерял %d HP",
				*diceDay1+*diceDay2,
				*diceDay1+*diceDay2,
			),
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
	}

	return dto.SleepyKingdomDTO{
		Exit:    false,
		Death:   death,
		NextTry: false,
	}, nil
}
