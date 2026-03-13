package bonus

import (
	"context"
	"fmt"
	"gamebook-backend/database/entities"
	battleCommon "gamebook-backend/modules/game/battle"
	"gamebook-backend/modules/game/bonus/helpers"
	"gamebook-backend/modules/game/dice"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"
	template2 "gamebook-backend/modules/game/template"

	"gorm.io/gorm"
)

var InstantHypnosisSpellName = "Заклинание Мгновенного Гипноза"
var InstantHypnosisSpellAlias = "instant_hypnosis_spell"

type instantHypnosisSpell struct {
	db     *gorm.DB
	player entities.Player
}

func NewInstantHypnosisSpell(db *gorm.DB, player entities.Player) dto.Bonus {
	return &instantHypnosisSpell{
		db:     db,
		player: player,
	}
}

func (h instantHypnosisSpell) Execute(ctx context.Context, req dto.BonusRequest) error {
	if len(h.player.Section.SectionEnemies) == 0 {
		return dto.ErrBattleNotFound
	}

	rollTheDices := dice.NewRollTheDices(h.db, &h.player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, h.player)
	if err != nil {
		return err
	}

	template, err := template2.GetDicesTemplate(ctx, *diceFirst, *diceSecond, false)
	if err != nil {
		return err
	}

	if *diceFirst+*diceSecond < 5 {
		err = h.Fail(ctx, fmt.Sprintf(
			"%s %s не сработало",
			template,
			InstantHypnosisSpellName,
		))
		if err != nil {
			return err
		}
	}

	description := fmt.Sprintf(
		"%s сработало, враг впадает в транс и больше вас не беспокоит",
		InstantHypnosisSpellName,
	)

	common, err := battleCommon.NewCommon(ctx, h.db, &h.player)
	if err != nil {
		return err
	}

	enemy := h.player.Section.SectionEnemies[0]
	oneEnemyHealth := enemy.Health

	health := uint(0)
	enemyUpdateListener, err := listener.HandleEvent(h.db, "enemy_update")
	if err != nil {
		return err
	}
	err = enemyUpdateListener.Handle(ctx, event.EnemyUpdateEvent{
		PlayerID:  h.player.ID,
		SectionID: h.player.SectionID,
		EnemyID:   enemy.ID,
		Health:    &health,
	})
	if err != nil {
		return err
	}

	battle := entities.Battle{
		Section:     h.player.Section.Number,
		PlayerID:    h.player.ID,
		Type:        dto.StepTypeSpell,
		Attacking:   battleCommon.AttackingPlayer,
		Damage:      oneEnemyHealth,
		Dice1:       *diceFirst,
		Dice2:       *diceSecond,
		Description: description,
		Weapon:      InstantHypnosisSpellAlias,
		Step:        common.Step().GetCurrentStepIndex(),
	}

	battleRepository := repository.NewBattleRepository(h.db)
	_, err = battleRepository.AddRecord(h.db, battle)
	if err != nil {
		return err
	}

	playerUpdateListener, err := listener.HandleEvent(h.db, "player_update")
	if err != nil {
		return err
	}

	bonusList := helpers.RemoveBonus(h.player.Bonus, InstantHypnosisSpellAlias)
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  h.player.ID,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(h.player.ID, fmt.Sprintf("<p>%s</p>", InstantHypnosisSpellName))

	return nil
}

func (h instantHypnosisSpell) Fail(ctx context.Context, description string) error {
	playerUpdateListener, err := listener.HandleEvent(h.db, "player_update")
	if err != nil {
		return err
	}
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID: h.player.ID,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(h.player.ID, fmt.Sprintf("<p>✨ %s</p>", description))

	return nil
}
