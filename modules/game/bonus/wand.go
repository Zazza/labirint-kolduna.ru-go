package bonus

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bonus/helpers"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/template"
	"math/rand"

	"gorm.io/gorm"
)

var WandName = "Волшебная палочка"
var WandAlias = "wand"

type wand struct {
	db     *gorm.DB
	player entities.Player
}

func NewWand(db *gorm.DB, player entities.Player) dto.Bonus {
	return &wand{
		db:     db,
		player: player,
	}
}

func (w wand) Execute(ctx context.Context, req dto.BonusRequest) error {
	if len(w.player.Section.SectionEnemies) == 0 {
		return dto.ErrBattleNotFound
	}

	enemyUpdateListener, err := listener.HandleEvent(w.db, "enemy_update")
	if err != nil {
		return err
	}

	playerUpdateListener, err := listener.HandleEvent(w.db, "player_update")
	if err != nil {
		return err
	}

	rollTheDices := dice.NewRollTheDices(w.db, &w.player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, w.player)
	if err != nil {
		return err
	}

	template2, err := template.GetDicesTemplate(ctx, *diceFirst, *diceSecond, false)
	if err != nil {
		return err
	}

	if *diceFirst+*diceSecond < 6 {
		err = w.Fail(ctx, fmt.Sprintf(
			"%s %s не сработала",
			template2,
			WandName,
		))
		if err != nil {
			return err
		}
	}

	duration := uint(4)
	debuff := entities.Debuff{
		Alias:    entities.DebuffAliasSkipReason,
		Duration: &duration,
	}

	randEnemy := rand.Intn(len(w.player.Section.SectionEnemies))
	enemy := w.player.Section.SectionEnemies[randEnemy]

	bonusList := helpers.RemoveBonus(w.player.Bonus, WandAlias)
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  w.player.ID,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(w.player.ID, fmt.Sprintf("<p>✨ %s</p>", WandName))

	err = enemyUpdateListener.Handle(ctx, event.EnemyUpdateEvent{
		PlayerID:    w.player.ID,
		SectionID:   w.player.SectionID,
		Debuff:      &[]entities.Debuff{debuff},
		EnemyID:     enemy.ID,
		Description: fmt.Sprintf("%s заморозила %s", WandName, enemy.Name),
	})
	if err != nil {
		return err
	}

	return nil
}

func (w wand) Fail(ctx context.Context, description string) error {
	playerUpdateListener, err := listener.HandleEvent(w.db, "player_update")
	if err != nil {
		return err
	}

	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID: w.player.ID,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(w.player.ID, fmt.Sprintf("<p>✨ %s</p>", description))

	return nil
}
