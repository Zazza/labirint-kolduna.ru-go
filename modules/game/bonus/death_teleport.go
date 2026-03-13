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

var DeathTeleportName = "Телепорт к месту гибели"
var DeathTeleportAlias = "death_teleport"

type deathTeleportSpell struct {
	db     *gorm.DB
	player entities.Player
}

func NewDeathTeleport(db *gorm.DB, player entities.Player) dto.Bonus {
	return &deathTeleportSpell{
		db:     db,
		player: player,
	}
}

func (d deathTeleportSpell) Execute(ctx context.Context, _ dto.BonusRequest) error {
	sectionRepository := repository.NewSectionRepository(d.db)
	playerSectionRepository := repository.NewPlayerSectionRepository(d.db)

	deathSection, err := sectionRepository.GetBySectionNumber(ctx, d.db, 9)
	if err != nil {
		return err
	}

	teleportSection, err := playerSectionRepository.GetPreviousSectionIdBySectionId(ctx, d.db, d.player.ID, deathSection.ID)
	playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
	if err != nil {
		return err
	}

	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  d.player.ID,
		SectionID: &teleportSection.SectionID,
	})
	if err != nil {
		return err
	}
	helper.DescriptionMessage(
		d.player.ID,
		fmt.Sprintf("%s переместил тебя в Секцию %d", DeathTeleportName, teleportSection.Section.Number),
	)

	bonusList := helpers.RemoveBonus(d.player.Bonus, DeathTeleportAlias)
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  d.player.ID,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	return nil
}
