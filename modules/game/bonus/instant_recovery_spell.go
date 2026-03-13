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

var InstantRecoveryName = "Заклинание Мгновенного Выздоровления"
var InstantRecoveryAlias = "instant_recovery_spell"

type instantRecoverySpell struct {
	db     *gorm.DB
	player entities.Player
}

func NewInstantRecoverySpell(db *gorm.DB, player entities.Player) dto.Bonus {
	return &instantRecoverySpell{
		db:     db,
		player: player,
	}
}

func (s instantRecoverySpell) Execute(ctx context.Context, req dto.BonusRequest) error {
	bonusList := helpers.RemoveBonus(s.player.Bonus, InstantRecoveryAlias)

	health := s.player.HealthMax
	playerUpdateListener, err := listener.HandleEvent(s.db, "player_update")
	if err != nil {
		return err
	}
	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  s.player.ID,
		Health:    &health,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(s.player.ID, fmt.Sprintf("<p>✨ %s</p>", InstantRecoveryName))

	return nil
}
