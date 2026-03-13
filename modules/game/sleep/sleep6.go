package sleep

import (
	"context"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"

	"gorm.io/gorm"
)

type sleep6 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep6(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep6{
		db:     db,
		player: player,
	}
}

func (s sleep6) Execute(
	ctx context.Context,
	dice1 uint,
	dice2 uint,
) (dto.SleepyKingdomDTO, error) {
	death := false
	exit := false
	if 6 <= dice1+dice2 {
		exit = true
	} else {
		rollTheDices := dice.NewRollTheDices(s.db, &s.player)

		dice1, dice2, err := rollTheDices.RollTheDices(ctx, s.player)
		if err != nil {
			return dto.SleepyKingdomDTO{}, err
		}

		if *dice1+*dice2 >= 6 {
			health := s.player.Health - 10

			helper.DescriptionMessage(
				s.player.ID,
				"<p>6-12 очков - попал на камни, потеряв 10 Жизненных Сил</p>",
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
		} else if 2 <= *dice1 && *dice2 < 6 {
			diceSwim1, diceSwim2, err := rollTheDices.RollTheDices(ctx, s.player)
			if err != nil {
				return dto.SleepyKingdomDTO{}, err
			}

			if *diceSwim1+*diceSwim2 >= 7 {
				exit = true

				helper.DescriptionMessage(
					s.player.ID,
					"<p>7-12 очков - выплыл благополучно</p>",
				)
			} else {
				death = true
			}
		}
	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   death,
		NextTry: false,
	}, nil
}
