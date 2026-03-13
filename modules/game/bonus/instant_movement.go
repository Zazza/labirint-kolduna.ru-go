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

var InstantMovementName = "Заклинание Мгновенного Перемещения"
var InstantMovementAlias = "instant_movement"

type instantMovement struct {
	db     *gorm.DB
	player entities.Player
}

func NewInstantMovement(db *gorm.DB, player entities.Player) dto.Bonus {
	return &instantMovement{
		db:     db,
		player: player,
	}
}

func (d instantMovement) Execute(ctx context.Context, req dto.BonusRequest) error {
	sectionsRepository := repository.NewSectionRepository(d.db)
	section, err := sectionsRepository.GetBySectionNumber(ctx, d.db, 3)
	if err != nil {
		return err
	}

	bonusList := helpers.RemoveBonus(d.player.Bonus, InstantMovementAlias)

	playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
	if err != nil {
		return err
	}
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  d.player.ID,
		SectionID: &section.ID,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", InstantMovementName))

	return nil
}
