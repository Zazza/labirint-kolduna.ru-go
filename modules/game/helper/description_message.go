package helper

import (
	"gamebook-backend/modules/game/channel"
	"gamebook-backend/modules/game/listener/event"

	"github.com/google/uuid"
)

func DescriptionMessage(playerID uuid.UUID, message string) {
	channel.ChPlayerSection <- event.PlayerSectionEvent{
		PlayerID:    playerID,
		Description: &message,
	}
}
