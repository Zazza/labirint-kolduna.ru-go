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

	"gorm.io/gorm"
)

var LuckyStoneName = "Счастливый Камешек"
var LuckyStoneAlias = "lucky_stone"

type luckyStone struct {
	db     *gorm.DB
	player entities.Player
}

func NewLuckyStone(db *gorm.DB, player entities.Player) dto.Bonus {
	return &luckyStone{
		db:     db,
		player: player,
	}
}

func (d luckyStone) Execute(ctx context.Context, req dto.BonusRequest) error {
	buff := entities.Buff{
		Alias: entities.DebuffAliasLuckyStoneReason,
	}

	bonusList := helpers.RemoveBonus(d.player.Bonus, LuckyStoneAlias)

	playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
	if err != nil {
		return err
	}
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  d.player.ID,
		Buff:      &[]entities.Buff{buff},
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", LuckyStoneName))

	return nil
}
