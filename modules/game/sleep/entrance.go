package sleep

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"
	template2 "gamebook-backend/modules/game/template"

	"gorm.io/gorm"
)

type Entrance interface {
	Handle(ctx context.Context) error
}

type entrance struct {
	db     *gorm.DB
	player entities.Player
}

func NewEntrance(db *gorm.DB, player entities.Player) Entrance {
	return &entrance{
		db:     db,
		player: player,
	}
}

func (e *entrance) Handle(ctx context.Context) error {
	sectionRepository := repository.NewSectionRepository(e.db)

	rollTheDices := dice.NewRollTheDices(e.db, &e.player)

	playerUpdateListener, err := listener.HandleEvent(e.db, "player_update")
	if err != nil {
		return err
	}

	diceSleepyKingdom, err := rollTheDices.RollTheDice(ctx, e.player)
	if err != nil {
		return err
	}

	templateSleepyKingdom, err := template2.GetDiceTemplate(ctx, *diceSleepyKingdom, false)
	if err != nil {
		return err
	}

	if *diceSleepyKingdom > 5 {
		diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, e.player)
		if err != nil {
			return err
		}

		templateHealthRecovery, err := template2.GetDicesTemplate(ctx, *diceFirst, *diceSecond, false)
		if err != nil {
			return err
		}

		healthRecovery := *diceFirst + *diceSecond
		helper.DescriptionMessage(
			e.player.ID,
			fmt.Sprintf(
				"<p>%s %s Вы хорошо поспали и восстановили %d HP</p>",
				templateSleepyKingdom,
				templateHealthRecovery,
				healthRecovery,
			),
		)

		var newPlayerHealth uint
		if e.player.HealthMax <= e.player.Health+healthRecovery {
			newPlayerHealth = e.player.HealthMax
		} else {
			newPlayerHealth = e.player.Health + healthRecovery
		}

		err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
			PlayerID: e.player.ID,
			Health:   &newPlayerHealth,
		})
		if err != nil {
			return err
		}

		return nil
	}

	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, e.player)
	if err != nil {
		return err
	}

	number := *diceFirst + *diceSecond

	goToSection, err := sectionRepository.GetBySectionNumber(ctx, e.db, 200+number)
	if err != nil {
		return err
	}

	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:        e.player.ID,
		SectionID:       &goToSection.ID,
		ReturnToSection: &e.player.SectionID,
	})
	if err != nil {
		return err
	}

	return nil
}
