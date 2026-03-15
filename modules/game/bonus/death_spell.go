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

var DeathSpellName = "Заклинание смерти"
var DeathSpellAlias = "death_spell"

type deathSpell struct {
	db     *gorm.DB
	player entities.Player
}

func NewDeathSpell(db *gorm.DB, player entities.Player) dto.Bonus {
	return &deathSpell{
		db:     db,
		player: player,
	}
}

func (d deathSpell) Execute(ctx context.Context, req dto.BonusRequest) error {
	common, err := battleCommon.NewCommon(ctx, d.db, &d.player)
	if err != nil {
		return err
	}

	enemies := common.GetEnemies()
	if enemies == nil || len(*enemies) == 0 {
		return dto.ErrBattleNotFound
	}

	// Фильтруем только живых врагов
	var aliveEnemies []entities.PlayerSectionEnemy
	for _, enemy := range *enemies {
		if enemy.Health > 0 {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}

	if len(aliveEnemies) == 0 {
		return dto.ErrBattleNotFound
	}

	// Бросаем кубики для проверки неудачи
	rollTheDices := dice.NewRollTheDices(d.db, &d.player)
	diceFirst, diceSecond, err := rollTheDices.RollTheDices(ctx, d.player)
	if err != nil {
		return err
	}

	template, err := template2.GetDicesTemplate(ctx, *diceFirst, *diceSecond, false)
	if err != nil {
		return err
	}

	// Проверка на неудачу: двойные 6, 1 или 3
	if (diceFirst == diceSecond && *diceFirst == 6) ||
		(diceFirst == diceSecond && *diceFirst == 1) ||
		(diceFirst == diceSecond && *diceFirst == 3) {
		err = d.Fail(ctx, fmt.Sprintf(
			"%s %s прикончило тебя",
			template,
			DeathSpellName,
		))
		if err != nil {
			return err
		}
	}

	// Выбираем случайного живого врага
	randomIndex := dice.RandomInt(len(aliveEnemies))
	enemy := aliveEnemies[randomIndex]
	oneEnemyHealth := enemy.Health

	health := uint(0)

	enemyUpdateListener, err := listener.HandleEvent(d.db, "enemy_update")
	if err != nil {
		return err
	}
	err = enemyUpdateListener.Handle(ctx, event.EnemyUpdateEvent{
		PlayerID:  d.player.ID,
		SectionID: d.player.SectionID,
		EnemyID:   enemy.EnemyID,
		Health:    &health,
	})
	if err != nil {
		return err
	}

	description := fmt.Sprintf("%s прикончило врага", DeathSpellName)

	battle := entities.Battle{
		Section:     d.player.Section.Number,
		PlayerID:    d.player.ID,
		Type:        dto.StepTypeSpell,
		Attacking:   battleCommon.AttackingPlayer,
		Damage:      oneEnemyHealth,
		Dice1:       *diceFirst,
		Dice2:       *diceSecond,
		Description: description,
		Weapon:      DeathSpellAlias,
		Step:        common.Step().GetCurrentStepIndex(),
	}

	battleRepository := repository.NewBattleRepository(d.db)
	_, err = battleRepository.AddRecord(d.db, battle)
	if err != nil {
		return err
	}

	playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
	if err != nil {
		return err
	}

	bonusList := helpers.RemoveBonus(d.player.Bonus, DeathSpellAlias)
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  d.player.ID,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", DeathSpellName))

	return nil
}

func (d deathSpell) Fail(ctx context.Context, description string) error {
	health := uint(0)

	playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
	if err != nil {
		return err
	}
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID: d.player.ID,
		Health:   &health,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", description))

	return nil
}
