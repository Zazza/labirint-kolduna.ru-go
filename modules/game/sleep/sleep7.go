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
	template2 "gamebook-backend/modules/game/template"

	"gorm.io/gorm"
)

type sleep7 struct {
	db     *gorm.DB
	player entities.Player
}

func NewSleep7(db *gorm.DB, player entities.Player) dto.Sleep {
	return &sleep7{
		db:     db,
		player: player,
	}
}

func (s *sleep7) Execute(
	ctx context.Context,
	dice1 uint,
	dice2 uint,
) (dto.SleepyKingdomDTO, error) {
	exit := false
	death := false
	if dice1+dice2 < 6 {
		exit = true

		helper.DescriptionMessage(
			s.player.ID,
			"<p>Набрал меньше 6 очков - ты жив и здоров</p>",
		)
	} else {
		rollTheDices := dice.NewRollTheDices(s.db, &s.player)

		dice1, dice2, err := rollTheDices.RollTheDices(ctx, s.player)
		if err != nil {
			return dto.SleepyKingdomDTO{}, err
		}

		template, err := template2.GetDicesTemplate(ctx, *dice1, *dice2, false)
		if err != nil {
			return dto.SleepyKingdomDTO{}, err
		}

		helper.SafeHTMLDescriptionMessage(
			s.player.ID,
			fmt.Sprintf("<p>%s</p>", template),
		)

		if s.player.Health < *dice1+*dice2 {
			death = true

			helper.DescriptionMessage(
				s.player.ID,
				"<p>Съел протухший кисель, заболел и умер</p>",
			)
		} else {
			exit = true
			health := s.player.Health - (*dice1 + *dice2)

			helper.DescriptionMessage(
				s.player.ID,
				fmt.Sprintf("<p>Съел протухший кисель заболел и потерял %d HP</p>", health),
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
	}

	return dto.SleepyKingdomDTO{
		Exit:    exit,
		Death:   death,
		NextTry: false,
	}, nil
}
