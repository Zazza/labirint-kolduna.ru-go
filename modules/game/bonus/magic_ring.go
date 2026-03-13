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

	"gorm.io/gorm"
)

var MagicRingName = "Магическое кольцо"
var MagicRingAlias = "magic_ring"

var MagicRingOptions = []string{"left", "right"}

type magicRing struct {
	db     *gorm.DB
	player entities.Player
}

func NewMagicRing(db *gorm.DB, player entities.Player) dto.Bonus {
	return &magicRing{
		db:     db,
		player: player,
	}
}

func (d magicRing) Execute(ctx context.Context, req dto.BonusRequest) error {
	playerUpdateListener, err := listener.HandleEvent(d.db, "player_update")
	if err != nil {
		return err
	}

	switch req.Option {
	case MagicRingOptions[0]:
		health := d.player.Health + uint(25)
		err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
			PlayerID: d.player.ID,
			Health:   &health,
		})
		if err != nil {
			return err
		}

		helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", MagicRingName))
	case MagicRingOptions[1]:
		for index, playerWeapon := range d.player.Weapons {
			if playerWeapon.Item == dto.Lightning {
				d.player.Weapons[index].Count++
			}
		}

		err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
			PlayerID: d.player.ID,
			Weapons:  &d.player.Weapons,
		})
		if err != nil {
			return err
		}

		helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", MagicRingName))
	default:
		return errors.New(fmt.Sprintf("%s invalid data", MagicRingAlias))
	}

	bonusList := helpers.RemoveBonus(d.player.Bonus, MagicRingAlias)

	err = playerUpdateListener.Handle(ctx, event.PlayerUpdateEvent{
		PlayerID:  d.player.ID,
		BonusList: &bonusList,
	})
	if err != nil {
		return err
	}

	helper.DescriptionMessage(d.player.ID, fmt.Sprintf("<p>✨ %s</p>", MagicRingName))

	return nil
}
