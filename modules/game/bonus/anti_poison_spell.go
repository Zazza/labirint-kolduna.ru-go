package bonus

import (
	"context"
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

var AntiPoisonSpellName = "Заклинание смерти"
var AntiPoisonSpellAlias = "anti_poison_spell"

type antiPoisonSpell struct {
	db     *gorm.DB
	player entities.Player
}

func NewAntiPoisonSpell(db *gorm.DB, player entities.Player) dto.Bonus {
	return &antiPoisonSpell{
		db:     db,
		player: player,
	}
}

func (d antiPoisonSpell) Execute(ctx context.Context, req dto.BonusRequest) error {
	for _, debuff := range d.player.Debuff {
		if debuff.Alias == entities.AliasPoisonReason {
			if d.player.Section.Number == dto.SectionDeath {
				playerSectionRepository := repository.NewPlayerSectionRepository(d.db)
				section, err := playerSectionRepository.GetLastPlayerSection(ctx, d.db, d.player.ID)
				if err != nil {
					return err
				}

				d.player.SectionID = section.SectionID
			}

			newHealth := d.player.Health + *debuff.Health

			debuffList := helpers.RemoveDebuff(d.player.Debuff, entities.AliasPoisonReason)
			bonusList := helpers.RemoveBonus(d.player.Bonus, AntiPoisonSpellAlias)

			playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
			if err != nil {
				return err
			}
			err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
				PlayerID:  d.player.ID,
				Health:    &newHealth,
				SectionID: &d.player.SectionID,
				Debuff:    &debuffList,
				BonusList: &bonusList,
			})
			if err != nil {
				return err
			}

			helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", AntiPoisonSpellName))
		}
	}

	return nil
}
