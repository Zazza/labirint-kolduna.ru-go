package bonus

import (
	"context"
	"errors"
	"fmt"
	"gamebook-backend/database/entities"
	"gamebook-backend/modules/game/bonus/helpers"
	"gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/helper"
	"gamebook-backend/modules/game/listener"
	"gamebook-backend/modules/game/listener/event"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

var MagicDuckName = "Магическая утка"
var MagicDuckAlias = "magic_duck"

var MagicDuckOptions = []string{"anti_magic", "section"}

type magicDuck struct {
	db     *gorm.DB
	player entities.Player
}

func NewMagicDuck(db *gorm.DB, player entities.Player) dto.Bonus {
	return &magicDuck{
		db:     db,
		player: player,
	}
}

func (d magicDuck) Execute(ctx context.Context, req dto.BonusRequest) error {
	switch req.Option {
	case MagicDuckOptions[0]:
		if len(d.player.Section.SectionEnemies) == 0 {
			return dto.ErrBattleNotFound
		}
		err := d.AntiMagic(ctx, d.player)
		if err != nil {
			return err
		}
	case MagicDuckOptions[1]:
		err := d.Section(ctx, d.player)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("%s invalid data", MagicDuckAlias))
	}

	return nil
}

func (d magicDuck) AntiMagic(ctx context.Context, player entities.Player) error {
	playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
	if err != nil {
		return err
	}

	enemyUpdateListener, err := listener.HandleEvent(d.db, "enemy_update")
	if err != nil {
		return err
	}

	for _, debuff := range player.Debuff {
		if debuff.Alias == entities.AliasMagicReason {
			player.Health += *debuff.Health

			debuffList := helpers.RemoveDebuff(d.player.Debuff, entities.AliasMagicReason)

			err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
				PlayerID: player.ID,
				Health:   debuff.Health,
				Debuff:   &debuffList,
			})
			if err != nil {
				return err
			}

			helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s восстановила магический ущерб</p>", AntiPoisonSpellName))
		}
	}

	debuff := entities.Debuff{
		Alias: entities.DebuffAliasMagicOffReason,
	}

	for _, enemy := range player.Section.SectionEnemies {
		err = enemyUpdateListener.Handle(ctx, event.EnemyUpdateEvent{
			PlayerID:    player.ID,
			SectionID:   player.SectionID,
			Debuff:      &[]entities.Debuff{debuff},
			EnemyID:     enemy.ID,
			Description: fmt.Sprintf("%s остановила магию у %s", MagicDuckName, enemy.Name),
		})
		if err != nil {
			return err
		}
	}

	bonusList := helpers.RemoveBonus(d.player.Bonus, MagicDuckAlias)

	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  d.player.ID,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", MagicDuckName))

	return nil
}

func (d magicDuck) Section(ctx context.Context, player entities.Player) error {
	playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
	if err != nil {
		return err
	}

	sectionsRepository := repository.NewSectionRepository(d.db)

	sectionNumber := uint(158)
	section, err := sectionsRepository.GetBySectionNumber(ctx, d.db, sectionNumber)
	if err != nil {
		return err
	}

	bonusList := helpers.RemoveBonus(d.player.Bonus, MagicDuckAlias)

	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  d.player.ID,
		SectionID: &section.ID,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", MagicDuckName))

	return nil
}
